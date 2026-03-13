package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rap_backend/config"
	"rap_backend/dao"
	"rap_backend/fileprocess"
	"rap_backend/rpc/bigdata"
	"rap_backend/utils"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
)

type GetSubTaskStaticsDTO struct {
	SubTaskId string `json:"subtask_id" form:"subtask_id"`
}

type GetSubTaskStaticsRetDTO struct {
	CallCount         int `json:"call_count" form:"call_count"`
	FinishedCallCount int `json:"finished_call_count" form:"finished_call_count"`
}

type CallInfoTranscription struct {
	Slots                     []interface{} `json:"slots"`
	Moment                    int           `json:"moment"`
	NewState                  string        `json:"newState"`
	Intention                 string        `json:"intention"`
	StateText                 string        `json:"stateText"`
	Transcript                string        `json:"transcript"`
	GLabelHitKeywords         string        `json:"gLabelHitKeywords"`
	InteractiveActionNameList []string      `json:"interactiveActionNameList"`
}

type EsCallInfo struct {
	Found  bool `json:"found"`
	Source struct {
		StoriesInfo []EsStoriesInfo `json:"stories_info"`
	} `json:"_source"`
}

type EsStoriesInfo struct {
	StoryName string `json:"story_name"`
	Intent    string `json:"intent"`
	States    []struct {
		AsrText   string `json:"asr_text"`
		StepName  string `json:"step_name"`
		Type      string `json:"type"`
		Condition struct {
			Func   string `json:"Func"`
			Input  string `json:"Input"`
			Expr   string `json:"Expr"`
			Result bool   `json:"Result"`
		} `json:"condition"`
		SlotVerifyList interface{} `json:"slot_verify_list"`
		StartTime      int64       `json:"start_time"`
		EndTime        int64       `json:"end_time"`
		NoVoiceTime    int         `json:"no_voice_time"`
		Actions        []struct {
			ActionName          string `json:"action_name"`
			TerminalAsrText     string `json:"terminal_asr_text"`
			AsrLastWordEndTime  int    `json:"asr_last_word_end_time"`
			IsVadSilence        int    `json:"is_vad_silence"`
			NextStep            string `json:"next_step"`
			NextStory           string `json:"next_story"`
			RobotStartSpeakTime int    `json:"robot_start_speak_time"`
			RobotEndSpeakTime   int    `json:"robot_end_speak_time"`
			ManStartSpeakTime   int    `json:"man_start_speak_time"`
			ManEndSpeakTime     int    `json:"man_end_speak_time"`
			ActionStart         int    `json:"action_start"`
			ActionEnd           int    `json:"action_end"`
			NoVoiceTime         int    `json:"no_voice_time"`
			InterruptTime       int    `json:"interrupt_time"`
			AsrWaitOverTime     int    `json:"asr_wait_over_time"`
			VadEndTime          int    `json:"vad_end_time"`
			Audios              []struct {
				AudioName string `json:"audio_name"`
				Text      string `json:"text"`
			} `json:"audios"`
			ResponseInfo []struct {
				AsrMsg string `json:"asr_msg"`
				Nlp    []struct {
					Type         string `json:"type"`
					Intent       string `json:"intent"`
					FixedIntent  string `json:"fixed_intent"`
					Confidence   int    `json:"confidence"`
					NlpStartTime int64  `json:"nlp_start_time"`
					NlpEndTime   int64  `json:"nlp_end_time"`
					Corpus       string `json:"corpus"`
				} `json:"nlp"`
				Slots []struct {
					Name   string `json:"name"`
					Result string `json:"result"`
				} `json:"slots"`
			} `json:"response_info"`
			AsrStatus         bool   `json:"AsrStatus"`
			VadStatus         bool   `json:"VadStatus"`
			AsrVadCheckStatus bool   `json:"AsrVadCheckStatus"`
			AsrFirstWord      string `json:"AsrFirstWord"`
		} `json:"actions"`
	} `json:"states"`
}

type TransStoryNlp struct {
	StepName string
	Nlp      string
	Slots    [][]string
}

type SentenceOrigin struct {
	Trans    interface{}
	TaskID   string
	CallID   string
	LabelIDs []string
}

func generateNewSubTaskId(taskId string) string {
	return taskId + "_sub_" + strconv.Itoa(time.Now().Nanosecond())
}

func CreateNewSubTask(taskInfo CreateTaskDTO, taskId string, operator string, offset int, subNum int) (error, string) {
	var subtaskDao dao.SubtaskInfo
	subtaskDao.StartTime, _ = time.Parse("2006-01-02 15:04:05", taskInfo.StartTime)
	subtaskDao.FinishTime, _ = time.Parse("2006-01-02 15:04:05", taskInfo.FinishTime)
	subtaskDao.SubtaskOperator = operator
	subtaskDao.SubtaskId = generateNewSubTaskId(taskId)
	subtaskDao.Status = utils.TASK_STATUS_NOTSTARTED
	err := subtaskDao.AddSubTaskInfo()
	for _, callid := range taskInfo.Callid[(offset):(offset + subNum)] {
		defaultValue, err4 := dao.GetLabelValueFromGaussInfo(callid)
		seelog.Infof("defaultvalue %s", defaultValue)
		if err4 != nil {
			seelog.Errorf("get label value from gaussinfo failed:%s", err.Error())
		}
		for _, labId := range taskInfo.LabelId {
			var labwork dao.LabelworkInfo
			labwork.SubtaskId = subtaskDao.SubtaskId
			labwork.LabelId = labId
			labInfo, e := dao.GetLabelInfoByLabelId(labId)
			if e != nil || labInfo == nil {
				seelog.Errorf("not exist label:", labId)
				continue
			}
			labwork.LabelName = labInfo.LabelName
			labwork.CallId = callid
			labwork.IsOptional = labInfo.IsOptional
			labwork.Status = utils.LABELWORK_STATUS_NOTSTART
			labwork.IsEditable = labInfo.IsEditable
			labwork.LabelDesc = labInfo.LabelDesc
			labinfo, err3 := dao.GetLabelInfoByLabelId(labId)
			if err3 != nil {
				seelog.Error("can't get this label %s", labinfo.LabelId)
				continue
			}
			if labinfo.IsEditable == 0 {
				if defaultValue != nil {
					if len(*defaultValue) > 0 {
						valueMap := (*defaultValue)[0]
						if _, ok := valueMap[labinfo.LabelName]; ok {
							seelog.Infof("label:%s, type:%s", labinfo.LabelName, reflect.TypeOf(valueMap[labinfo.LabelName]).String())
							switch infoValue := valueMap[labinfo.LabelName].(type) {
							case string:
								labwork.LabelValue = infoValue
							case float32:
								labwork.LabelValue = fmt.Sprintf("%.2f", infoValue)
							case float64:
								labwork.LabelValue = fmt.Sprintf("%.2f", infoValue)
							case int:
							case int16:
							case int32:
							case int64:
								labwork.LabelValue = fmt.Sprintf("%d", infoValue)
							case time.Time:
								labwork.LabelValue = infoValue.Format("2006-01-02 15:04:05")
							default:
								labwork.LabelValue = fmt.Sprintf("%s", infoValue)
							}
						}
					}
				}
			}
			err2 := labwork.AddLabelworkInfo()
			if err2 != nil {
				seelog.Errorf("add label work info failed:", err2.Error())
				continue
			}
		}
	}

	return err, subtaskDao.SubtaskId
}

func GetSubTaskStatics(subTaskId string) *GetSubTaskStaticsRetDTO {
	result, err := dao.GetAllCallidListBySubtaskId(subTaskId)
	if err != nil {
		seelog.Errorf("add label work info failed:", err.Error())
		return nil
	}
	var ret GetSubTaskStaticsRetDTO
	ret.CallCount = len(result)
	finishedCnt := 0
	for _, v := range result {
		if v == utils.LABELWORK_STATUS_FINISHED {
			finishedCnt++
		}
	}
	ret.FinishedCallCount = finishedCnt
	return &ret
}

func GetSubTaskDownloadUrl(subtaskId string) (string, error) {
	var excelTable [][]string
	labworklst, err2 := dao.GetAllLabelWorkListBySubtaskId(subtaskId)
	if err2 != nil {
		seelog.Errorf("get label work infos failed:%s", err2.Error())
		return "", err2
	}
	if labworklst == nil {
		seelog.Errorf("no content for subtask:%s", subtaskId)
		return "", errors.New("no content for subtask:" + subtaskId)
	}
	allLabelWork := make(map[string]map[string]string)
	for _, labwork := range *labworklst {
		if _, ok := allLabelWork[labwork.CallId]; !ok {
			allLabelWork[labwork.CallId] = make(map[string]string)
		}
		allLabelWork[labwork.CallId][labwork.LabelName] = labwork.LabelValue
	}
	headerList := []string{
		"call_id",
	}
	firstFlg := 1
	for callid, mapvalue := range allLabelWork {
		var rowValue []string
		rowValue = append(rowValue, callid)
		for k, v := range mapvalue {
			rowValue = append(rowValue, v)
			if firstFlg == 1 {
				headerList = append(headerList, k)
			}
		}
		firstFlg = 0
		excelTable = append(excelTable, rowValue)
	}

	seelog.Infof("excelTable:%s", excelTable)

	xlserr := utils.OutPutDataWithXLSX2(excelTable, headerList, fileprocess.CFG_LOCALREPORTFILEPATH, subtaskId+".xlsx")
	if xlserr != nil {
		seelog.Errorf("generate xls file failed %s", xlserr.Error())
		return "", xlserr
	}
	return utils.EXCEL_DOWNLOAD_PATH + subtaskId + ".xlsx", nil
}

// 添加 新建任务-callids
func CreateTaskCalls(taskID string, callIDs []string) error {
	if len(callIDs) == 0 {
		return nil
	}
	tcs := make([]dao.TaskCall, 0)
	for k, callid := range callIDs {
		tc := dao.TaskCall{
			TaskId:       taskID,
			CallId:       callid,
			Status:       utils.TASK_STATUS_CREATED,
			SerialNumber: k + 1,
		}
		tcs = append(tcs, tc)
	}
	err := dao.CreateTaskCallInBatches(tcs)
	return err
}

// 添加 新建任务-callid下-labels info
func CreateTaskLabel(taskInfo CreateTaskDTO, taskId string) error {
	labDatas := make([]dao.LabelworkInfo, 0)
	sentenseDatas := make([]SentenceOrigin, 0)
	for _, callid := range taskInfo.Callid {
		defaultValue, err := dao.GetLabelValueFromGaussInfo(callid)
		seelog.Infof("defaultvalue %s", defaultValue)
		if err != nil {
			seelog.Errorf("get label value from gaussinfo failed:%s", err.Error())
		}
		for _, labId := range taskInfo.LabelId {
			var labwork dao.LabelworkInfo
			labwork.SubtaskId = taskId
			labwork.TaskId = taskId
			labwork.LabelId = labId
			labInfo, ok := LabelInfoCache[labId]
			if !ok {
				seelog.Errorf("not exist label:%s", labId)
				continue
			}
			labwork.LabelName = labInfo.LabelName
			labwork.CallId = callid
			labwork.IsOptional = labInfo.LabelIsOptional
			labwork.Status = utils.TASK_STATUS_CREATED
			labwork.IsEditable = labInfo.LabelIsEditable
			labwork.LabelDesc = labInfo.LabelDesc
			labwork.LabelTime = utils.NowUTC()
			if labInfo.LabelIsEditable == 0 {
				if defaultValue != nil {
					if len(*defaultValue) > 0 {
						valueMap := (*defaultValue)[0]
						if _, ok := valueMap[labInfo.LabelName]; ok {
							seelog.Infof("label:%s, type:%s", labInfo.LabelName, reflect.TypeOf(valueMap[labInfo.LabelName]).String())
							switch infoValue := valueMap[labInfo.LabelName].(type) {
							case string:
								labwork.LabelValue = infoValue
							case float32:
								labwork.LabelValue = fmt.Sprintf("%.2f", infoValue)
							case float64:
								labwork.LabelValue = fmt.Sprintf("%.2f", infoValue)
							case int:
							case int16:
							case int32:
							case int64:
								labwork.LabelValue = fmt.Sprintf("%d", infoValue)
							case time.Time:
								labwork.LabelValue = infoValue.Format("2006-01-02 15:04:05")
							default:
								labwork.LabelValue = fmt.Sprintf("%s", infoValue)
							}
						} else if labInfo.LabelName == "Sentence" {
							//获取Transcription
							if val, ok := valueMap["Transcription"]; ok {
								sentence := "" //TranscriptionToSentence(val, callid)
								senOri := SentenceOrigin{
									Trans:    val,
									TaskID:   taskId,
									CallID:   callid,
									LabelIDs: []string{labInfo.LabelId},
								}
								labwork.LabelValue = sentence
								if labInfo1, ok := LabelNameInfoCache["Sentence1"]; ok {
									labSen1 := dao.LabelworkInfo{
										SubtaskId:  taskId,
										TaskId:     taskId,
										LabelId:    labInfo1.LabelId,
										LabelName:  labInfo1.LabelName,
										LabelValue: sentence,
										CallId:     callid,
										IsOptional: labInfo1.LabelIsOptional,
										Status:     utils.TASK_STATUS_CREATED,
										IsEditable: labInfo1.LabelIsEditable,
										LabelDesc:  labInfo1.LabelDesc,
										LabelTime:  utils.NowUTC(),
									}
									labDatas = append(labDatas, labSen1)
									senOri.LabelIDs = append(senOri.LabelIDs, labInfo1.LabelId)
								}
								sentenseDatas = append(sentenseDatas, senOri)

							}
						} else if labInfo.LabelName == "RobotType" {
							if val, ok := valueMap["RobotName"]; ok {
								switch infoValue := val.(type) {
								case string:
									labwork.LabelValue = getRobotTypeByRobotName(infoValue)
								}
							}
						}
					}
				}
			}
			labDatas = append(labDatas, labwork)
		}
	}
	err := dao.CreateTaskLabelWorkInBatches(labDatas)
	if err != nil {
		seelog.Errorf("add label work info failed:%s", err.Error())
		return err
	}
	if len(sentenseDatas) > 0 {
		go taskSentenseDataInit(sentenseDatas)
	}

	return nil
}

// 原来的函数从gausscallinfo中获取数据，现在从大数据接口获取数据
func CreateTaskLabel2(taskInfo CreateTaskDTO, taskId string) error {
	go createTaskLabel(taskInfo, taskId)
	return nil
}

func concatSentence(sentences []string) (ret string) {
	for index, sentence := range sentences {
		sentence = strings.Replace(sentence, "<", "{", -1)
		sentence = strings.Replace(sentence, ">", "}", -1)
		sentences[index] = sentence
	}

	return strings.Join(sentences, "\n")
}

func createTaskLabel(taskInfo CreateTaskDTO, taskId string) {
	seelog.Infof("start handler create task label, taskId: %s", taskId)
	defer seelog.Infof("end handler create task label, taskId: %s", taskId)

	callInfosMap, err := bigdata.GetSvc().GetCallInfoMap(bigdata.GetCallInfoReq{
		CallIds: taskInfo.Callid,
	})
	if err != nil {
		seelog.Errorf("failed to get call info, taskId: %s, err: %+v", taskId, err)
		callInfosMap = make(map[string]*bigdata.CallInfo)
		for _, callId := range taskInfo.Callid {
			callInfo, err := bigdata.GetSvc().GetCallInfo(bigdata.GetCallInfoReq{CallIds: []string{callId}})
			if err != nil {
				seelog.Errorf("failed to get call info, taskId: %s, callId: %s, err: %+v", taskId, callId, err)
				continue
			} else {
				callInfosMap[callId] = callInfo
			}
		}
	}

	for _, callId := range taskInfo.Callid {

		labDatas := make([]dao.LabelworkInfo, 0)

		if _, ok := callInfosMap[callId]; !ok {
			seelog.Errorf("failed to get callInfo from bigdata, callId: %s", callId)
			continue
		}
		callInfo := callInfosMap[callId]

		for _, labId := range taskInfo.LabelId {
			var labwork dao.LabelworkInfo
			labwork.SubtaskId = taskId
			labwork.TaskId = taskId
			labwork.LabelId = labId
			labInfo, ok := LabelInfoCache[labId]
			if !ok {
				seelog.Errorf("not exist taskId: %s, label:%s", taskId, labId)
				continue
			}
			labwork.LabelName = labInfo.LabelName
			labwork.CallId = callId
			labwork.IsOptional = labInfo.LabelIsOptional
			labwork.Status = utils.TASK_STATUS_CREATED
			labwork.IsEditable = labInfo.LabelIsEditable
			labwork.LabelDesc = labInfo.LabelDesc
			labwork.LabelTime = utils.NowUTC()

			switch labInfo.LabelName {
			case "CallID":
				labwork.LabelValue = callInfo.CallId
			case "RobotName":
				labwork.LabelValue = callInfo.RobotName
			case "StartTime":
				labwork.LabelValue = callInfo.StartTime
			case "Intention":
				labwork.LabelValue = callInfo.Intention
			case "BillSec":
				labwork.LabelValue = fmt.Sprintf("%d", callInfo.BillSec)
			case "Sentence":
				labwork.LabelValue = concatSentence(callInfo.SenTence)
				if labInfo1, ok := LabelNameInfoCache["Sentence1"]; ok {
					labSen1 := dao.LabelworkInfo{
						SubtaskId:  taskId,
						TaskId:     taskId,
						LabelId:    labInfo1.LabelId,
						LabelName:  labInfo1.LabelName,
						LabelValue: concatSentence(callInfo.SenTence),
						CallId:     callId,
						IsOptional: labInfo1.LabelIsOptional,
						Status:     utils.TASK_STATUS_CREATED,
						IsEditable: labInfo1.LabelIsEditable,
						LabelDesc:  labInfo1.LabelDesc,
						LabelTime:  utils.NowUTC(),
					}
					labDatas = append(labDatas, labSen1)
				}
			case "TalkRound":
				labwork.LabelValue = fmt.Sprintf("%d", callInfo.TalkRound)
			case "RobotType":
				labwork.LabelValue = callInfo.RobotType
			case "CalleeNumber":
				labwork.LabelValue = callInfo.CalleeNumber
			case "SIPLine":
				labwork.LabelValue = callInfo.SIPLine
			case "Company":
				labwork.LabelValue = callInfo.Company
			}

			labDatas = append(labDatas, labwork)
		}

		err := dao.CreateTaskLabelWorkInBatches(labDatas)
		if err != nil {
			seelog.Errorf("add label work info, taskId: %s, callId: %s, failed:%s", taskId, callId, err.Error())
			continue
		}
	}

	return
}

// 添加 ringtype 新建任务-callid下-labels info
func CreateRingTypeTaskLabel(callids []RingTypeCallInfo, labelIDs []string, taskId, systemLabelID, countryLabelID string) error {
	labDatas := make([]dao.LabelworkInfo, 0)
	for _, c := range callids {
		callid := c.CallID
		for _, labId := range labelIDs {
			var labwork dao.LabelworkInfo
			labwork.SubtaskId = taskId
			labwork.TaskId = taskId
			labwork.LabelId = labId
			labInfo, ok := LabelInfoCache[labId]
			if !ok {
				seelog.Errorf("not exist label:%s", labId)
				continue
			}
			labwork.LabelName = labInfo.LabelName
			labwork.CallId = callid
			labwork.IsOptional = labInfo.LabelIsOptional
			labwork.Status = utils.TASK_STATUS_CREATED
			labwork.IsEditable = labInfo.LabelIsEditable
			labwork.LabelDesc = labInfo.LabelDesc
			labwork.LabelTime = utils.NowUTC()
			if labId == systemLabelID {
				labwork.LabelValue = c.SystemRingType
			}
			if labId == countryLabelID {
				labwork.LabelValue = c.Country
			}
			labDatas = append(labDatas, labwork)
		}

	}
	err := dao.CreateTaskLabelWorkInBatches(labDatas)
	if err != nil {
		seelog.Errorf("add label work info failed:%s", err.Error())
		return err
	}
	return nil
}

// ringtype 录音文件地址入库
func PrepareRingTypeTaskRecordingFiles(taskId string, callidList []RingTypeCallInfo) {
	tcs := make([]dao.CallidRecording, 0)
	for _, c := range callidList {
		tc := dao.CallidRecording{
			CallId:  c.CallID,
			Address: c.URL,
		}
		tcs = append(tcs, tc)
	}
	err := dao.CreateCallRecordingInBatches(tcs)
	if err != nil {
		seelog.Errorf("CreateCallRecordingInBatches failed:%s", err.Error())
		return
	}
	err5 := dao.UpdateTaskStatusByTaskID(taskId, utils.TASK_STATUS_CREATED)
	if err5 != nil {
		seelog.Errorf("PrepareRingTypeTaskRecordingFiles task %s status %s failed:%s", taskId, utils.TASK_STATUS_CREATED, err5.Error())
		return
	}
}

// gauss_callinfo中的Transcription 转换成 Sentence想要的数据及格式
func TranscriptionToSentence(trans interface{}, callID string) string {
	text := ""
	switch infoValue := trans.(type) {
	case string:
		txt := make([]CallInfoTranscription, 0)
		err := json.Unmarshal([]byte(infoValue), &txt)
		if err != nil {
			seelog.Errorf("transcriptionToSentence Unmarshal callid:%s, txt:%s, err:%s", callID, infoValue, err.Error())
			return ""
		}
		text = TranSliceToSentence(txt, callID)
	default:
		return text
	}
	return text
}

func TranSliceToSentence(trans []CallInfoTranscription, callID string) string {
	storyInfo := GetStoryInfo(callID)
	texts := []string{}
	for k, tran := range trans {
		text := tran.Transcript + " [" + tran.Intention + "]"
		if len(storyInfo) > k && tran.NewState == storyInfo[k].StepName {
			story := storyInfo[k]
			sl := []string{}
			for _, slot := range story.Slots {
				sl = append(sl, strings.Join(slot, ":"))
			}
			if len(sl) > 0 {
				text += " [" + strings.Join(sl, " ") + "]"
			}
			if story.Nlp != "" {
				text += " [STORY_NLU: " + story.Nlp + "]"
			}
		}
		texts = append(texts, text)
	}
	return strings.Join(texts, "\n")
}

func GetStoryInfo(callID string) []TransStoryNlp {
	data := make([]TransStoryNlp, 0)
	resp := EsCallInfo{}
	url := strings.Join([]string{config.EsDomain, config.EsIndex, config.EsDoc, callID}, "/")
	err := utils.DoHttpGetJson(url, &resp, http.Header{})
	if err != nil {
		seelog.Errorf("GetStoryInfo DoHttpGetJson err:%s", err.Error())
		return data
	}
	if !resp.Found || len(resp.Source.StoriesInfo) == 0 {
		return data
	}
	for _, stories := range resp.Source.StoriesInfo {
		for _, state := range stories.States {
			story := TransStoryNlp{
				StepName: state.StepName,
			}
			for _, action := range state.Actions {
				for _, response := range action.ResponseInfo {
					for _, nlp := range response.Nlp {
						if nlp.Type == "STORY_NLU" {
							story.Nlp = nlp.Intent
						}
					}
					story.Slots = make([][]string, 0)
					slotSlice := []string{}
					for _, slot := range response.Slots {
						if slot.Name == "未填上" {
							slotSlice = append(slotSlice, slot.Name)
							continue
						}
						slotSlice = append(slotSlice, slot.Name)
						slotSlice = append(slotSlice, slot.Result)
					}
					if len(slotSlice) > 0 {
						story.Slots = append(story.Slots, slotSlice)
					}
				}
			}
			data = append(data, story)
		}
	}
	return data
}

func taskSentenseDataInit(datas []SentenceOrigin) {
	for _, sent := range datas {
		sentence := TranscriptionToSentence(sent.Trans, sent.CallID)
		if sentence == "" {
			continue
		}
		upd := map[string]interface{}{
			"label_value": sentence,
		}
		_, err := dao.UpdateLabelworkInfoByCallIDLabID(sent.TaskID, sent.CallID, sent.LabelIDs, upd)
		if err != nil {
			seelog.Errorf("UpdateLabelworkInfoByCallIDLabID err:%s, taskID:%s, callid:%s, sent:%s", err.Error(), sent.TaskID, sent.CallID, sentence)
		}
	}
}

func getRobotTypeByRobotName(robotName string) string {
	robotType := ""
	if strings.HasPrefix(robotName, "PO") || strings.HasPrefix(robotName, "Promotion") {
		robotType = "promotion"
	} else if strings.HasPrefix(robotName, "M0") || strings.HasPrefix(robotName, "M1") {
		robotType = "collection"
	} else if strings.HasPrefix(robotName, "IC") {
		robotType = "Infocheck"

	}
	return robotType
}
