package service

import (
	"fmt"
	"os"
	"path"
	"rap_backend/dao"
	"rap_backend/fileprocess"
	"rap_backend/rpc/bigdata"
	"rap_backend/rpc/dmrobot"
	"rap_backend/utils"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

var PTE_SHOW_LABELS = map[string]bool{
	"RobotType": true,
	"RobotName": true,
	"StartTime": true,
	"CallID":    true,
	"BillSec":   true,
	"Intention": true,
	"TalkRound": true,
	"Sentence":  true,
	"RingType":  true,
}

type GetOneCallLabelWorkDTO struct {
	CallId    string `json:"call_id" form:"call_id"`
	SubTaskId string `json:"subtask_id" form:"subtask_id"`
	TaskId    string `json:"task_id" form:"task_id"`
}

type LabelWorkInfoDTO struct {
	LabelInfoItem
	AuditorContent string `json:"auditor_content" form:"auditor_content"`
	LKID           int32  `json:"lk_id"`
	LabelValue     string `json:"label_value" form:"label_value"`
	LabelStatus    string `json:"label_status" form:"label_status"`
}

type ExcelLabelWorkInfo struct {
	LabelValueList map[string]string `json:"label_value_list" form:"label_value_list"`
}

type UpdateOneCallLabelWorkDTO struct {
	CallId           string             `json:"call_id" form:"call_id"`
	TaskId           string             `json:"task_id" form:"task_id"`
	LabelWorkInfoLst []LabelWorkInfoDTO `json:"label_work_info_lst" form:"label_work_info_lst"`
}

type UpdateRejectOneCallLabelWorkDTO struct {
	RejectReason     string             `json:"reject_reason" form:"reject_reason"`
	CallId           string             `json:"call_id" form:"call_id"`
	TaskId           string             `json:"task_id" form:"task_id"`
	IsAll            bool               `json:"is_all" form:"is_all"`
	LabelWorkInfoLst []LabelWorkInfoDTO `json:"label_work_info_lst" form:"label_work_info_lst"`
}

type GetOneCallLabelWorkRetDTO struct {
	CallId string `json:"call_id"`
	CallRecordURL
	RejectReason     string             `json:"reject_reason" form:"reject_reason"`
	LabelWorkInfoLst []LabelWorkInfoDTO `json:"label_work_info_lst"`
	Creator          string             `json:"creator"`
	Annotator        string             `json:"annotator"`
	Auditor          string             `json:"auditor"`
}

type CallRecordURL struct {
	RecordUrl    string `json:"record_url" form:"record_url"`
	RecordUrlCh1 string `json:"record_url_ch1" form:"record_url_ch1"`
	RecordUrlCh2 string `json:"record_url_ch2" form:"record_url_ch2"`
}

type GetSubtaskLabelworkListDTO struct {
	SubtaskIdList []string `json:"subtask_id_list" form:"subtask_id_list"`
	TaskID        string   `json:"task_id"`
	PageNum       int      `json:"page_num" form:"page_num"`
	PageSize      int      `json:"page_size" form:"page_size"`
}

type GetSubtaskLabelworkListRetDTO struct {
	TotalCallidNum  int                           `json:"total_num" form:"total_num"`
	LabelworkDetail map[string][]LabelWorkInfoDTO `json:"label_work_info" form:"label_work_info_lst"`
}
type TaskPreviewLabelworkListDTO struct {
	TaskID   string `json:"task_id"`
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type TaskPreviewLabelworkListRetDTO struct {
	TaskTotal int                           `json:"task_total" form:"task_total"`
	Total     int                           `json:"total" form:"total"`
	List      map[string][]LabelWorkInfoDTO `json:"list" form:"list"`
	Progress  []*dao.TaskProgressDetail     `json:"progress" form:"progress"`
}

type TaskCallLabelListItem struct {
	CallID string             `json:"call_id"`
	Labels []LabelWorkInfoDTO `json:"labels" form:"labels"`
}

type TaskLabelWorkDoneDTO struct {
	TaskId string `json:"task_id" form:"task_id"`
	Total  int    `json:"total" form:"total"`
}

type TaskCallLabelListDown struct {
	CallID string              `json:"call_id"`
	Labels []LabelWorkInfoDown `json:"labels" form:"labels"`
}

type LabelWorkInfoDown struct {
	LabelId        string `json:"label_id" form:"label_id"`
	LabelName      string `json:"label_name" form:"label_name"`
	IsColor        int    `json:"is_color" form:"is_color"`
	AuditorContent string `json:"auditor_content" form:"auditor_content"`
	LabelValue     string `json:"label_value" form:"label_value"`
}

func GetSubtaskLabelworkList(subtaskId []string, pageNum int, pageSize int) (*GetSubtaskLabelworkListRetDTO, error) {
	labelworkList, err := dao.GetLabelWorkListBySubtaskId2(subtaskId, pageNum, pageSize)
	callidTotal := 0
	if err != nil {
		seelog.Errorf("GetSubtaskLabelworkList failedL:%s", err.Error())
		return nil, err
	} else {
		var ret GetSubtaskLabelworkListRetDTO
		labelWorkListMap := make(map[string][]LabelWorkInfoDTO)
		for _, labelworkDetail := range *labelworkList {
			if labelWorkListMap[labelworkDetail.CallId] == nil {
				labelWorkListMap[labelworkDetail.CallId] = []LabelWorkInfoDTO{}
			}
			var detail LabelWorkInfoDTO
			detail.LabelId = labelworkDetail.LabelId
			detail.LabelName = labelworkDetail.LabelName
			detail.LabelValue = labelworkDetail.LabelValue
			detail.LabelStatus = labelworkDetail.Status
			detail.LabelIsOptional = labelworkDetail.IsOptional
			detail.LabelIsEditable = labelworkDetail.IsEditable
			detail.LabelDesc = labelworkDetail.LabelDesc
			labelWorkListMap[labelworkDetail.CallId] = append(labelWorkListMap[labelworkDetail.CallId], detail)
		}
		ret.LabelworkDetail = labelWorkListMap
		for _, sid := range subtaskId {
			callidList, err3 := dao.GetAllCallidListBySubtaskId(sid)
			if err3 != nil {
				seelog.Errorf("get callid list failed, subtask id:%s", sid)
				continue
			}
			callidTotal = callidTotal + len(callidList)
		}
		ret.TotalCallidNum = callidTotal
		return &ret, nil
	}
}

// 获取callid下 所有label works列表
func GetOneCallLabelWork(taskid, callid, pageTab string) []LabelWorkInfoDTO {
	var detail0 = make([]LabelWorkInfoDTO, 0)
	var detail1 = make([]LabelWorkInfoDTO, 0)
	var detail2 = make([]LabelWorkInfoDTO, 0) //Sentence1 Chinese 排在同一行
	var detail3 = make([]LabelWorkInfoDTO, 0) //Sentence1 Chinese 排在同一行
	var details = make([]LabelWorkInfoDTO, 0)

	detail, err := dao.GetLabelWorkDetail(taskid, callid)
	if err != nil {
		seelog.Errorf("get label work detail failed:%s", err)
		return details
	}

	for _, value := range *detail {
		lab, ok := LabelInfoCache[value.LabelId]
		if !ok {
			continue
		}
		labelworkVal := LabelWorkInfoDTO{
			LabelInfoItem: LabelInfoItem{
				LabelId:         lab.LabelId,
				LabelName:       lab.LabelName,
				LabelIsOptional: lab.LabelIsOptional,
				LabelIsEditable: lab.LabelIsEditable,
				LabelDesc:       lab.LabelDesc,
				LabelType:       lab.LabelType,
				LabelOptions:    lab.LabelOptions,
				IsColor:         lab.IsColor,
			},
			LKID:           value.Id,
			LabelValue:     value.LabelValue,
			LabelStatus:    value.Status,
			AuditorContent: value.AuditorContent,
		}
		if value.AuditorContent != "" {
			labelworkVal.LabelValue = value.AuditorContent
		}
		switch pageTab {

		//rap3.0 ，结果分析页面放开编辑
		case utils.TASK_PAGE_TAB_ANNOTAT, utils.TASK_PAGE_TAB_AUDIT, utils.TASK_PAGE_TAB_ANALYS: //标注、审核页面时，结果分析字段不展示
			//belong 含义：1只在结果分析页面展示并编辑，2标注、审核、分析页面都可展示并编辑
			if lab.Belong == 1 {
				continue
			}
			if lab.LabelIsEditable == 0 {
				if _, ok := PTE_SHOW_LABELS[lab.LabelName]; !ok {
					continue
				}
			}
			//case utils.TASK_PAGE_TAB_ANALYS: //结果分析页面，其他字段不可编辑
			//	if lab.Belong != 1 && lab.Belong != 2 {
			//		labelworkVal.LabelIsEditable = 0
			//	}
		}
		if lab.LabelName == "Sentence1" {
			detail2 = append(detail2, labelworkVal)
		} else if lab.LabelName == "Chinese" {
			detail3 = append(detail3, labelworkVal)
		} else if lab.LabelIsEditable == 0 { //不可编辑 靠前展示
			detail0 = append(detail0, labelworkVal)
		} else {
			detail1 = append(detail1, labelworkVal)
		}
	}
	details = append(details, detail0...)
	details = append(details, detail2...)
	details = append(details, detail3...)
	details = append(details, detail1...)

	//特殊处理 real intention字段
	for index, labelInfo := range details {
		if labelInfo.LabelName == "Real Intention" {
			callInfo, err := bigdata.GetSvc().GetCallInfo(bigdata.GetCallInfoReq{
				CallIds: []string{callid},
			})
			if err != nil {
				continue
			}
			annotateList, err := dmrobot.GetSvc().GetRobotIntentionAnnotateList(dmrobot.GetRobotIntentionAnnotateListReq{
				RobotName: callInfo.RobotName,
			})
			if err != nil {
				continue
			}

			labelDesc := labelInfo.LabelDesc
			for _, annotate := range annotateList {
				labelDesc = fmt.Sprintf("%s\n%s-%s-%s", labelDesc, annotate.Intention, annotate.EnAnnotate, annotate.ZhAnnotate)
			}

			details[index].LabelDesc = labelDesc
		}
	}

	return details
}

func GetRecordURLByTaskIDCallID(taskID, callID string) CallRecordURL {
	var cr = CallRecordURL{}
	var rcids = []string{callID, callID + "_a", callID + "_b"}
	subPath := path.Join(fileprocess.CFG_LOCALRECORDINGFILEPATH, taskID) + "/"

	for _, cid := range rcids {
		isExist, _ := PathExists(subPath + cid + ".wav")
		if isExist {
			switch cid {
			case callID:
				cr.RecordUrl = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+".wav")
			case callID + "_a":
				cr.RecordUrlCh1 = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+"_a.wav")
			case callID + "_b":
				cr.RecordUrlCh2 = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+"_b.wav")

			}
		}
	}
	if cr.RecordUrl != "" {
		return cr
	}

	//先去callid recording表获取一下
	recs, err := dao.GetCallRecordingInfoByCallID(rcids)
	if err == nil && recs != nil && len(*recs) > 0 {
		for _, rec := range *recs {
			if rec.CallId == callID {
				cr.RecordUrl = fileprocess.CFG_RECORDINGURL + rec.Address
			}
			if rec.CallId == callID+"_a" {
				cr.RecordUrlCh1 = fileprocess.CFG_RECORDINGURL + rec.Address
			}
			if rec.CallId == callID+"_b" {
				cr.RecordUrlCh2 = fileprocess.CFG_RECORDINGURL + rec.Address
			}
		}
		return cr
	}
	cr.RecordUrl = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+".wav")
	cr.RecordUrlCh1 = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+"_a.wav")
	cr.RecordUrlCh2 = fileprocess.CFG_RECORDINGURL + path.Join(taskID, callID+"_b.wav")
	return cr
}

// callid 保存标注内容
func UpdateOneCallLabelWork(updateLabelWorkInfoLst []LabelWorkInfoDTO) error {
	for _, content := range updateLabelWorkInfoLst {
		daoLable := dao.LabelworkInfo{
			Id:         content.LKID,
			LabelTime:  time.Now(),
			LabelValue: content.LabelValue,
			Status:     utils.TASK_STATUS_ANNOTATED,
		}
		columns := []string{"status"}
		if content.LabelIsEditable == 1 {
			columns = []string{"status", "label_value"}
		}
		_, err := daoLable.UpdateLabelworkInfo(columns)
		if err != nil {
			seelog.Errorf("update label work detail failed:%s", err)
		}
	}
	return nil
}

// callid 保存审核内容
func UpdateAuditOneCallLabelWork(updateLabelWorkInfoLst []LabelWorkInfoDTO) error {
	for _, content := range updateLabelWorkInfoLst {
		daoLable := dao.LabelworkInfo{
			Id:             content.LKID,
			LabelTime:      time.Now(),
			AuditorContent: content.LabelValue,
			Status:         utils.TASK_STATUS_AUDITED,
		}
		columns := []string{"status"}
		if content.LabelIsEditable == 1 {
			columns = []string{"status", "auditor_content"}
		}
		_, err := daoLable.UpdateLabelworkInfo(columns)
		if err != nil {
			seelog.Errorf("update audit label work detail failed:%s", err)
		}
	}
	return nil
}

// callid 保存结果分析内容
func UpdateAnalystOneCallLabelWork(updateLabelWorkInfoLst []LabelWorkInfoDTO) error {
	for _, content := range updateLabelWorkInfoLst {
		daoLable := dao.LabelworkInfo{
			Id:             content.LKID,
			LabelTime:      time.Now(),
			AuditorContent: content.LabelValue,
			Status:         utils.TASK_STATUS_ANALYSTED,
		}
		columns := []string{"status"}
		if content.LabelIsEditable == 1 {
			columns = []string{"status", "auditor_content"}
		}
		_, err := daoLable.UpdateLabelworkInfo(columns)
		if err != nil {
			seelog.Errorf("update audit label work detail failed:%s", err)
		}
	}
	return nil
}

// 获取处理callid的用户 创建者 标注者 审核者
func GetOneCallLabelWorkUsers(taskID, callID string) (string, string, string) {
	tc, err := dao.GetTaskCallInfoByTaskIDCallID(taskID, callID)
	if err != nil || tc == nil {
		return "", "", ""
	}
	task, err := dao.GetTaskInfoByID(taskID)
	if err != nil || tc == nil {
		return "", "", ""
	}
	actualCreator := task.TaskCreator
	actualAnnotator := ""
	actualAuditor := ""
	uids := []uint32{}
	var annuid uint32
	var auduid uint32
	if tc.ActualAnnotator != "" {
		uid, _ := strconv.ParseInt(tc.ActualAnnotator, 10, 64)
		if uid > 0 {
			annuid = uint32(uid)
			uids = append(uids, uint32(uid))
		}
	}
	if tc.ActualAuditor != "" {
		uid, _ := strconv.ParseInt(tc.ActualAuditor, 10, 64)
		if uid > 0 {
			auduid = uint32(uid)
			uids = append(uids, uint32(uid))
		}
	}
	if len(uids) > 0 {
		users, _ := UserShortInfoByIDs(uids)
		if u, ok := users[annuid]; ok {
			actualAnnotator = u.UserName
		}
		if u, ok := users[auduid]; ok {
			actualAuditor = u.UserName
		}
	}
	return actualCreator, actualAnnotator, actualAuditor
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
