package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"rap_backend/dao"
	"rap_backend/fileprocess"
	"rap_backend/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/golibs/uuid"
)

type DownloadTaskDTO struct {
	TaskId    string `json:"task_id" form:"task_id"`
	SubTaskId string `json:"subtask_id" form:"subtask_id"`
}

type TaskDetailDTO struct {
	TaskId string `json:"task_id" form:"task_id" binding:"required"`
}

type SubmitStatisticsDTO struct {
	TaskId string `json:"task_id" form:"task_id" binding:"required"`
	Type   string `json:"type" form:"type" binding:"required"`
}

type DownloadTaskRetDTO struct {
	DownloadUrl string `json:"download_url" form:"download_url"`
}

type CheckCallIdDTO struct {
	Callid []string `json:"call_id" form:"call_id"`
}

type CheckCallIdRetDTO struct {
	CheckedCallid   []string `json:"checked_call_id" form:"checked_call_id"`
	ExceptionCallid []string `json:"exception_call_id" form:"exception_call_id"`
}

type LabelDefaultRetDTO struct {
	Labels []string `json:"labels" form:"labels"`
}

type CreateTaskDTO struct {
	Callid        []string `json:"call_id" form:"call_id"`
	TaskCreator   string   `json:"task_creator" form:"task_creator"`
	CreatorId     uint32
	TaskName      string    `json:"task_name" form:"task_name"`
	TaskOperators []uint32  `json:"task_operators" form:"task_operators"`
	StartTime     string    `json:"start_time" form:"start_time"`
	FinishTime    string    `json:"finish_time" form:"finish_time"`
	LabelId       []string  `json:"label_id" form:"label_id"`
	StartTimer    time.Time `json:"-"`
	FinishTimer   time.Time `json:"-"`
}

type CreateTaskRetDTO struct {
	Taskid string `json:"task_id" form:"task_id"`
}

type TaskInfo struct {
	Id             int32    `json:"-"`
	Taskid         string   `json:"task_id" form:"task_id"`
	TaskName       string   `json:"task_name" form:"task_name"`
	CallNum        int      `json:"call_num" form:"call_num"`
	TaskCreator    string   `json:"task_creator" form:"task_creator"`
	StartTime      string   `json:"start_time" form:"start_time"`
	FinishTime     string   `json:"finish_time" form:"finish_time"`
	TaskCreateTime string   `json:"task_create_time" form:"task_create_time"`
	TaskStatus     string   `json:"task_status" form:"task_status"`
	TaskOperator   []string `json:"task_operator" form:"task_operator"`
	SubTask        []string `json:"sub_task" form:"sub_task"`
	Allocator      []uint32 `json:"-"`
	Annotator      []uint32 `json:"-"`
	Auditor        []uint32 `json:"-"`
	Analysts       []uint32 `json:"-"`
}

type GetTaskListDTO struct {
	SearchType  string `json:"search_type" form:"search_type"`
	SearchValue string `json:"search_value" form:"search_value"`
	PageNum     int    `json:"page_num" form:"page_num"`
	PageSize    int    `json:"page_size" form:"page_size"`
}

type GetTaskListRetDTO struct {
	TaskCnt      int64      `json:"task_cnt" form:"task_cnt"`
	TaskInfoList []TaskInfo `json:"task_info_list" form:"task_info_list"`
}

type GetTaskCallIdListDTO struct {
	TaskId    string `json:"task_id" form:"task_id"`
	SubTaskId string `json:"subtask_id" form:"subtask_id"`
	PageNum   int    `json:"page_num" form:"page_num"`
	PageSize  int    `json:"page_size" form:"page_size"`
	NoMark    int    `json:"nomark" form:"nomark"`
}

type CallIdStatus struct {
	CallId       string `json:"call_id" form:"call_id"`
	Status       string `json:"status" form:"status"`
	SerialNumber int    `json:"serial_number" from:"serial_number"`
}

type GetTaskCallIdListRetDTO struct {
	TaskId       string           `json:"task_id" form:"task_id"`
	SubtaskId    string           `json:"subtask_id" form:"subtask_id"`
	CallIdList   []CallIdStatus   `json:"call_id_list" form:"call_id_list"`
	Total        int64            `json:"total"`
	LastLocation TaskLastLocation `json:"last_location"`
	Statistics   struct {
		Total int64 `json:"total"`
		Done  int64 `json:"done"`
	} `json:"statistics"`
}

type TaskLastLocation struct {
	Page int `json:"page_num"`
	CallIdStatus
}

type TaskListFilter struct {
	PageTab     string
	StatusList  []string
	UserIDs     []uint32
	SearchType  string
	SearchValue string
	PageNum     int
	PageSize    int
	TIDs        []int32
	CurrentUID  uint32
}

type DoAllocatTaskDTO struct {
	Id             int32                  `json:"-"`
	TaskId         string                 `json:"task_id" form:"task_id" binding:"required"`
	Annotator      []dao.DoTaskUserNum    `json:"annotator" form:"annotator"`
	Auditor        []dao.DoTaskUserNum    `json:"auditor" form:"auditor"`
	Analysts       []dao.DoTaskUserNum    `json:"analysts" form:"analysts"`
	Type           string                 `json:"type" form:"type"`
	TaskAnnotator  string                 `json:"-"`
	TaskAuditor    string                 `json:"-"`
	TaskAnalysts   string                 `json:"-"`
	Relation       []TaskUserRelationItem `json:"-"`
	StatusInit     string                 `json:"-"`
	TaskStatusInit string                 `json:"-"`
	NumAnnot       map[uint32]int         `json:"-"`
	NumAudit       map[uint32]int         `json:"-"`
	NumAnal        map[uint32]int         `json:"-"`
}

type RedoAllocatTaskSingleDTO struct {
	Id           int32                  `json:"-"`
	TaskId       string                 `json:"task_id" form:"task_id" binding:"required"`
	Operator     []dao.DoTaskUserNum    `json:"operator" form:"operator"`
	Type         string                 `json:"type" form:"type"`
	TaskOperator string                 `json:"-"`
	Relation     []TaskUserRelationItem `json:"-"`
	NumAudit     map[uint32]int         `json:"-"`
}

type TaskUserRelationItem struct {
	TurID    int32
	TID      int32
	TaskID   string
	UserID   uint32
	TaskRole string
	TotalNum int
	DoneNum  int
	Status   string
}

type GetTaskStatDTO struct {
	TaskId string `json:"task_id" form:"task_id"`
}

type GetTaskStatRetDTO struct {
	Total    int `json:"total" form:"total"`
	Finished int `json:"finished" form:"finished"`
}

type RingTypeTaskDTO struct {
	Key        string             `json:"key"`
	TaskName   string             `json:"task_name"`
	StartTime  string             `json:"start_time"`
	FinishTime string             `json:"finish_time"`
	CallIDs    []RingTypeCallInfo `json:"callids"`
	LableNames []string           `json:"label_names"`
}

type RingTypeCallInfo struct {
	CallID         string `json:"call_id"`
	SystemRingType string `json:"system_ring_type"`
	Country        string `json:"country"`
	URL            string `json:"url"`
}

type TaskCallItem struct {
	TcId            int32  `json:"tc_id" from:"tc_id"`
	TaskId          string `json:"task_id" from:"task_id"`
	CallId          string `json:"call_id" form:"call_id"`
	MultiAnnotator  uint32 `json:"multi_annotator" form:"multi_annotator"`
	MultiAuditor    uint32 `json:"multi_auditor" form:"multi_auditor"`
	MultiAnalysts   uint32 `json:"multi_analysts" form:"multi_analysts"`
	ActualAnnotator uint32 `json:"actual_annotator" form:"actual_annotator"`
	ActualAuditor   uint32 `json:"actual_auditor" form:"actual_auditor"`
	ActualAnalysts  uint32 `json:"actual_analysts" form:"actual_analysts"`
	SerialNumber    int    `json:"serial_number" from:"serial_number"`
	Status          string `json:"status" form:"status"`
	RejectReason    string `json:"reject_reason" form:"reject_reason"`
}

func generateNewTaskId() string {
	return "tsk_" + strings.Replace(uuid.Rand().Hex(), "-", "", -1)
}

func CreateNewTask(taskInfo CreateTaskDTO) (error, string) {
	var taskDao dao.TaskInfo
	taskDao.TaskName = taskInfo.TaskName
	taskDao.TaskCreatetime = taskInfo.StartTimer
	taskDao.TaskCreator = taskInfo.TaskCreator
	taskDao.CreatorId = int(taskInfo.CreatorId)
	taskDao.CallNum = len(taskInfo.Callid)
	taskDao.MultiAllocator = utils.UInt32SliceToString(taskInfo.TaskOperators)
	taskDao.Labels = strings.Join(taskInfo.LabelId, ",")

	seelog.Infof("starttime:%s", taskInfo.StartTime)
	taskDao.StartTime = taskInfo.StartTimer
	seelog.Infof("daostarttime:%s", taskDao.StartTime)
	taskDao.FinishTime = taskInfo.FinishTimer

	taskDao.Status = utils.TASK_STATUS_CREATED
	taskDao.TaskId = generateNewTaskId()
	err := taskDao.AddTaskInfo()
	return err, taskDao.TaskId
}

// func GetTaskList(searchType string, searchValue string, pageNum int, pageSize int) (*[]TaskInfo, int64, error) {
// 	seelog.Infof("GetTaskList, pagenum:%d, pagesize:%d", pageNum, pageSize)
// 	totalCnt, err0 := dao.GetTaskListCount()
// 	if err0 != nil {
// 		seelog.Errorf("get task list count failed:%s", err0.Error())
// 		return nil, 0, err0
// 	}
// 	tasklist, err := dao.GetTaskInfoList(searchType, searchValue, pageNum, pageSize)
// 	if err != nil {
// 		seelog.Errorf("get task list failed:%s", err.Error())
// 		return nil, 0, err
// 	}
// 	if tasklist == nil {
// 		return nil, 0, nil
// 	}
// 	var retTaskList = make([]TaskInfo, 0)
// 	for _, taskinfo := range *tasklist {
// 		var info TaskInfo
// 		info.Taskid = taskinfo.TaskId
// 		info.TaskName = taskinfo.TaskName
// 		info.StartTime = taskinfo.StartTime.Format("2006-01-02 15:04:05")
// 		info.FinishTime = taskinfo.FinishTime.Format("2006-01-02 15:04:05")
// 		info.TaskCreator = taskinfo.TaskCreator
// 		info.TaskCreateTime = taskinfo.TaskCreatetime.Format("2006-01-02 15:04:05")
// 		info.TaskStatus = taskinfo.Status
// 		info.CallNum = taskinfo.CallNum
// 		distribution, err2 := dao.GetDistributedSubTaskByTaskid(taskinfo.TaskId)
// 		if err2 != nil || distribution == nil {
// 			seelog.Errorf("get task distribute failed:%s", err.Error())
// 			return nil, 0, err2
// 		}
// 		optflg := 0
// 		for _, dis := range *distribution {
// 			subTaskInfo, err3 := dao.GetSubTaskInfoBySubTaskid(dis.SubtaskId)
// 			if err3 != nil {
// 				seelog.Errorf("get task distribute failed:%s", err3.Error())
// 				continue
// 			}
// 			if searchType == utils.TASK_SEARCH_TYPE_OPERATOR {
// 				if strings.Contains(strings.ToLower(subTaskInfo.SubtaskOperator), strings.ToLower(searchValue)) {
// 					optflg = 1
// 					info.TaskOperator = append(info.TaskOperator, subTaskInfo.SubtaskOperator)
// 					info.SubTask = append(info.SubTask, subTaskInfo.SubtaskId)
// 				}
// 			} else {
// 				info.TaskOperator = append(info.TaskOperator, subTaskInfo.SubtaskOperator)
// 				info.SubTask = append(info.SubTask, subTaskInfo.SubtaskId)
// 			}

// 		}
// 		if searchType == utils.TASK_SEARCH_TYPE_OPERATOR {
// 			if optflg == 1 {
// 				retTaskList = append(retTaskList, info)
// 			}
// 		} else {
// 			retTaskList = append(retTaskList, info)
// 		}
// 	}
// 	return &retTaskList, totalCnt, nil
// }

func CheckUploadCallId(checkCallIdLst []string) ([]string, []string, error) {
	extlst, notextlst, err := dao.CheckTaskCallId2(checkCallIdLst)
	if err != nil {
		seelog.Errorf("check callid failed:%s", err.Error())
		return nil, nil, err
	}
	return extlst, notextlst, err
}

func GetTaskCallIdList(taskId string, subtaskId string, pageNum int, pageSize int) (*GetTaskCallIdListRetDTO, error) {
	if subtaskId == "" {
		taskDistribute, err := dao.GetDistributedSubTaskByTaskid(taskId)
		if err != nil || taskDistribute == nil {
			seelog.Error("Get subTask id failed or no task distribute")
			return nil, err
		}
		subtaskId = (*taskDistribute)[0].SubtaskId
		seelog.Infof("get distribute subtaskid:%s", subtaskId)
	}
	labelWorkLst, err := dao.GetLabelWorkListBySubtaskId(subtaskId, pageNum, pageSize)
	if err != nil || labelWorkLst == nil {
		seelog.Error("Get Task CallId List failed or no label work list")
		return nil, err
	}
	var taskCallIdList GetTaskCallIdListRetDTO
	var callIdList []CallIdStatus
	for _, labelwork := range *labelWorkLst {
		var callstatus CallIdStatus
		callstatus.CallId = labelwork.CallId
		callstatus.Status = labelwork.Status
		callIdList = append(callIdList, callstatus)
	}
	taskCallIdList.CallIdList = callIdList
	taskCallIdList.SubtaskId = subtaskId
	taskCallIdList.TaskId = taskId
	return &taskCallIdList, nil
}

func GetTaskDownloadUrl(taskId string) (string, error) {
	taskDistribute, err := dao.GetDistributedSubTaskByTaskid(taskId)
	if err != nil {
		seelog.Errorf("get %s distribute faild:%s", taskId, err)
		return "", err
	}
	if taskDistribute == nil {
		seelog.Errorf("no content for task:%s", taskId)
		return "", errors.New("no content for task:" + taskId)
	}

	var excelTable [][]string
	headerList := []string{
		"call_id",
	}
	for _, dis := range *taskDistribute {
		labworklst, err2 := dao.GetAllLabelWorkListBySubtaskId(dis.SubtaskId)
		if err2 != nil {
			seelog.Errorf("get label work infos failed:%s", err2.Error())
			return "", err2
		}
		if labworklst == nil {
			seelog.Errorf("no content for subtask:%s", dis.SubtaskId)
			return "", errors.New("no content for subtask:" + dis.SubtaskId)
		}
		allLabelWork := make(map[string]map[string]string)
		for _, labwork := range *labworklst {
			if _, ok := allLabelWork[labwork.CallId]; !ok {
				allLabelWork[labwork.CallId] = make(map[string]string)
			}
			allLabelWork[labwork.CallId][labwork.LabelName] = labwork.LabelValue
		}

		var tmpCallid = make([]string, 0)
		for k, _ := range allLabelWork {
			tmpCallid = append(tmpCallid, k)
		}
		sort.Strings(tmpCallid)
		firstFlg := 1
		for _, callid := range tmpCallid {
			mapvalue := allLabelWork[callid]
			var rowValue []string
			rowValue = append(rowValue, callid)
			var tmpKeys = make([]string, 0)
			for k, _ := range mapvalue {
				tmpKeys = append(tmpKeys, k)
			}
			sort.Strings(tmpKeys)
			for _, value := range tmpKeys {
				if firstFlg == 1 {
					headerList = append(headerList, value)
				}
				rowValue = append(rowValue, mapvalue[value])
			}
			firstFlg = 0
			excelTable = append(excelTable, rowValue)
		}
	}

	seelog.Infof("excelTable:%s", excelTable)

	xlserr := utils.OutPutDataWithXLSX2(excelTable, headerList, fileprocess.CFG_LOCALREPORTFILEPATH, taskId+".xlsx")
	if xlserr != nil {
		seelog.Errorf("generate xls file failed %s", xlserr.Error())
		return "", xlserr
	}

	return utils.EXCEL_DOWNLOAD_PATH + taskId + ".xlsx", nil
}

// func DeleteTask(taskId string) error {
// 	taskDistribute, err := dao.GetDistributedSubTaskByTaskid(taskId)
// 	if err != nil {
// 		seelog.Errorf("get %s distribute faild:%s", taskId, err)
// 		return err
// 	}
// 	if taskDistribute == nil {
// 		seelog.Errorf("no content for task:%s", taskId)
// 		return errors.New("no content for task:" + taskId)
// 	}
// 	for _, dis := range *taskDistribute {
// 		err2 := dao.DeleteSubtaskLabelWork(dis.SubtaskId)
// 		if err2 != nil {
// 			seelog.Errorf("DeleteSubtaskLabelWork %s  faild:%s", dis.SubtaskId, err2)
// 			return err2
// 		}
// 		err3 := dao.DeleteSubtaskInSubTaskInfos(dis.SubtaskId)
// 		if err3 != nil {
// 			seelog.Errorf("DeleteSubtaskInSubTaskInfos %s  faild:%s", dis.SubtaskId, err3)
// 			return err3
// 		}
// 	}
// 	err4 := dao.DeleteTaskInDistribute(taskId)
// 	if err4 != nil {
// 		seelog.Errorf("DeleteTaskInDistribute %s  faild:%s", taskId, err4)
// 		return err4
// 	}
// 	err5 := dao.DeleteTaskInTaskInfos(taskId)
// 	if err5 != nil {
// 		seelog.Errorf("DeleteTaskInTaskInfos %s  faild:%s", taskId, err5)
// 		return err5
// 	}
// 	return nil
// }

// 更新任务状态.
func UpdateTaskStatusByTaskID(taskId, status string) error {
	err := dao.UpdateTaskStatusByTaskID(taskId, status)
	if err == nil && (status == utils.TASK_STATUS_COMPLETED || status == utils.TASK_STATUS_DELETED) {
		//清除本地音频文件
		go ClearTaskLocalWavFile(taskId)
	}
	return err
}

// 更新任务状态only
func UpdateTaskStatusByTaskIDOnly(taskId, status string) error {
	return dao.UpdateTaskStatusByTaskID(taskId, status)
}

func CheckTaskName(taskName string) (string, error) {
	task, err := dao.CheckTaskNameExist(taskName)
	if err != nil {
		return "", err
	}
	if task != nil { //taskname已存在
		return task.TaskId, fmt.Errorf("任务已存在")
	}
	return "", nil
}

// 筛选任务列表
func GetTaskListByFilter(req TaskListFilter) (*[]TaskInfo, int64, error) {
	seelog.Infof("GetTaskList, pagenum:%d, pagesize:%d", req.PageNum, req.PageSize)
	reqDao := dao.TaskListFilter{
		PageTab:     req.PageTab,
		StatusList:  req.StatusList,
		UserIDs:     req.UserIDs,
		SearchType:  req.SearchType,
		SearchValue: req.SearchValue,
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
	}
	//个人用户 做已完成的任务筛选
	if len(req.UserIDs) == 1 {
		tabStatus := ""
		switch req.PageTab {
		case utils.TASK_PAGE_TAB_ANNOTAT:
			tabStatus = utils.TASK_STATUS_ANNOTATED
		case utils.TASK_PAGE_TAB_AUDIT:
			tabStatus = utils.TASK_STATUS_AUDITED
		case utils.TASK_PAGE_TAB_ANALYS:
			tabStatus = utils.TASK_STATUS_ANNOTATED
		}
		if tabStatus != "" {
			tids, err := dao.GetTaskIDsByUserRole(req.PageTab, tabStatus, req.UserIDs[0])
			if err == nil && len(tids) > 0 {
				reqDao.TIDs = tids
			}
		}
	}
	tasklist, totalCnt, err := dao.GetTaskListByFilter(reqDao)
	if err != nil {
		seelog.Errorf("get task list failed:%s", err.Error())
		return nil, 0, err
	}
	if tasklist == nil || len(*tasklist) == 0 {
		return nil, totalCnt, nil
	}
	var (
		tids            = []int32{}
		userUnqMap      = make(map[uint32]bool, 0)
		userUnqs        = []uint32{}
		userRelationMap = make(map[int32]int, 0)
	)
	for _, taskinfo := range *tasklist {
		tids = append(tids, taskinfo.Id)
		//当前负责人
		priIDs := []uint32{}
		switch taskinfo.Status {
		case utils.TASK_STATUS_CREATED:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAllocator)
		case utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATING:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAnnotator)
		case utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_AUDITING:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAuditor)
		case utils.TASK_STATUS_AUDITED, utils.TASK_STATUS_ANALYSTING, utils.TASK_STATUS_ANALYSTED, utils.TASK_STATUS_COMPLETED:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAnalysts)
		}
		for _, pid := range priIDs {
			if _, ok := userUnqMap[pid]; ok {
				continue
			}
			userUnqMap[pid] = true
			userUnqs = append(userUnqs, pid)
		}
	}

	userMap, err := UserShortInfoByIDs(userUnqs)
	if err != nil {
		seelog.Errorf("UserShortInfoByIDs err:%s, uids:%v", err.Error(), userUnqs)
	}

	//个人用户 做已完成的任务筛选
	if len(req.UserIDs) == 1 {
		userRelation, err := dao.GetTaskUserRelationByUIDRole(req.UserIDs[0], req.PageTab, tids)
		if err == nil && userRelation != nil && len(*userRelation) > 0 {
			for _, ur := range *userRelation {
				userRelationMap[ur.Tid] = ur.TotalNum
			}
		}
	}

	country := GetCountryInfoByUID(req.CurrentUID)
	var retTaskList = make([]TaskInfo, 0)
	for _, taskinfo := range *tasklist {
		var info TaskInfo
		info.Taskid = taskinfo.TaskId
		info.TaskName = taskinfo.TaskName
		info.StartTime = utils.UTCTime2TimeZone(taskinfo.StartTime, country.ZoneId)
		info.FinishTime = utils.UTCTime2TimeZone(taskinfo.FinishTime, country.ZoneId)
		info.TaskCreator = taskinfo.TaskCreator
		info.TaskCreateTime = utils.UTCTime2TimeZone(taskinfo.TaskCreatetime, country.ZoneId)
		info.TaskStatus = taskinfo.Status
		info.TaskOperator = []string{}
		info.SubTask = []string{}
		info.CallNum = taskinfo.CallNum
		if cn, ok := userRelationMap[taskinfo.Id]; ok && cn > 0 {
			info.CallNum = cn
		}
		//当前负责人
		priIDs := []uint32{}
		switch taskinfo.Status {
		case utils.TASK_STATUS_CREATED:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAllocator)
		case utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATING:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAnnotator)
		case utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_AUDITING:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAuditor)
		case utils.TASK_STATUS_AUDITED, utils.TASK_STATUS_ANALYSTING, utils.TASK_STATUS_ANALYSTED, utils.TASK_STATUS_COMPLETED:
			priIDs = utils.StringToUInt32Slice(taskinfo.MultiAnalysts)
		}
		for _, pid := range priIDs {
			if u, ok := userMap[pid]; ok {
				info.TaskOperator = append(info.TaskOperator, u.UserName)
			}
		}
		info.Allocator = utils.StringToUInt32Slice(taskinfo.MultiAllocator)
		info.Annotator = utils.StringToUInt32Slice(taskinfo.MultiAnnotator)
		info.Auditor = utils.StringToUInt32Slice(taskinfo.MultiAuditor)
		info.Analysts = utils.StringToUInt32Slice(taskinfo.MultiAnalysts)
		retTaskList = append(retTaskList, info)

	}
	return &retTaskList, totalCnt, nil
}

// 根据taskid 获取taskinfo.
func GetTaskInfoByID(taskID string) (*TaskInfo, error) {
	taskinfo, err := dao.GetTaskInfoByID(taskID)
	if err != nil {
		seelog.Errorf("GetTaskInfoByID failed:%s", err.Error())
		return nil, err
	}
	if taskinfo == nil {
		return nil, nil
	}
	var info TaskInfo
	info.Id = taskinfo.Id
	info.Taskid = taskinfo.TaskId
	info.TaskName = taskinfo.TaskName
	info.StartTime = taskinfo.StartTime.Format("2006-01-02 15:04:05")
	info.FinishTime = taskinfo.FinishTime.Format("2006-01-02 15:04:05")
	info.TaskCreator = taskinfo.TaskCreator
	info.TaskCreateTime = taskinfo.TaskCreatetime.Format("2006-01-02 15:04:05")
	info.TaskStatus = taskinfo.Status
	info.TaskOperator = []string{}
	info.SubTask = []string{}
	info.CallNum = taskinfo.CallNum
	info.Allocator = utils.StringToUInt32Slice(taskinfo.MultiAllocator)
	info.Annotator = utils.StringToUInt32Slice(taskinfo.MultiAnnotator)
	info.Auditor = utils.StringToUInt32Slice(taskinfo.MultiAuditor)
	info.Analysts = utils.StringToUInt32Slice(taskinfo.MultiAnalysts)
	return &info, nil
}

func GetOriginTaskInfoByID(taskID string, needAllocatInfo bool) (*dao.TaskAndAllocatInfo, error) {
	info, err := dao.GetTaskInfoByID(taskID)
	if err != nil {
		return nil, err
	}
	if !needAllocatInfo {
		return &dao.TaskAndAllocatInfo{
			TaskInfo: info,
		}, nil
	}

	//批量插入是有序的，根据id正序计算
	taskuserrelationList, err := dao.GetTaskUserRelationListByTaskId(taskID)
	if err != nil {
		return nil, err
	}

	Annotator := make([]dao.DoTaskUserNum, 0)
	Auditor := make([]dao.DoTaskUserNum, 0)
	Analysts := make([]dao.DoTaskUserNum, 0)

	for _, relation := range taskuserrelationList {
		if relation.TaskRole == utils.TASK_PAGE_TAB_AUDIT {
			Auditor = append(Auditor, dao.DoTaskUserNum{
				UserID: relation.UserId,
				Num:    relation.TotalNum,
			})
		} else if relation.TaskRole == utils.TASK_PAGE_TAB_ANNOTAT {
			Annotator = append(Annotator, dao.DoTaskUserNum{
				UserID: relation.UserId,
				Num:    relation.TotalNum,
			})

		} else if relation.TaskRole == utils.TASK_PAGE_TAB_ANALYS {
			Analysts = append(Analysts, dao.DoTaskUserNum{
				UserID: relation.UserId,
				Num:    relation.TotalNum,
			})
		}
	}

	return &dao.TaskAndAllocatInfo{
		TaskInfo:  info,
		Annotator: Annotator,
		Auditor:   Auditor,
		Analysts:  Analysts,
	}, nil
}

// 分配任务
func DoAllocatTask(req DoAllocatTaskDTO, mode string) error {
	//ids
	ids, err := dao.GetTcIDsByTaskID(req.TaskId)
	if err != nil {
		return err
	}
	initStatus := utils.TASK_STATUS_ALLOCATED
	if req.StatusInit != "" {
		initStatus = req.StatusInit
	}
	//task calls表更新
	k := 0
	k1 := 0
	for _, a1 := range req.Annotator {
		if a1.Num == 0 {
			continue
		}
		k1 += a1.Num
		upd := map[string]interface{}{
			"multi_annotator": a1.UserID,
			"status":          initStatus,
		}
		_, err := dao.UpdateTaskCallByTcIDs(ids[k:k1], upd)
		if err != nil {
			seelog.Errorf("UpdateTaskCallByTcIDs Annotator data:%+v, err:%s", a1, err.Error())
		}
		k = k1
	}
	k = 0
	k1 = 0
	for _, a1 := range req.Auditor {
		if a1.Num == 0 {
			continue
		}
		k1 += a1.Num
		upd := map[string]interface{}{
			"multi_auditor": a1.UserID,
			"status":        initStatus,
		}
		_, err := dao.UpdateTaskCallByTcIDs(ids[k:k1], upd)
		if err != nil {
			seelog.Errorf("UpdateTaskCallByTcIDs Operator data:%+v, err:%s", a1, err.Error())
		}
		k = k1
	}
	k = 0
	k1 = 0
	for _, a1 := range req.Analysts {
		if a1.Num == 0 {
			continue
		}
		k1 += a1.Num
		upd := map[string]interface{}{
			"multi_analysts": a1.UserID,
			"status":         initStatus,
		}
		_, err := dao.UpdateTaskCallByTcIDs(ids[k:k1], upd)
		if err != nil {
			seelog.Errorf("UpdateTaskCallByTcIDs Analysts data:%+v, err:%s", a1, err.Error())
		}
		k = k1
	}

	tk := dao.TaskInfo{
		Id:             req.Id,
		MultiAnnotator: req.TaskAnnotator,
		MultiAuditor:   req.TaskAuditor,
		MultiAnalysts:  req.TaskAnalysts,
		//Status:         initStatus,
	}

	if mode != "redo" {
		tk.Status = req.TaskStatusInit
	}

	_, err = tk.UpdateTaskInfoByID("multi_annotator", "multi_auditor", "multi_analysts", "status")

	//redo 需要删除历史task_user_relation
	if mode == "redo" {
		err = dao.DeleteTaskUserRelationByTaskID(req.TaskId)
		if err != nil {
			return err
		}
	}

	//任务 用户 角色 关系表 入库
	rels := []dao.TaskUserRelation{}
	for _, r := range req.Relation {
		rel := dao.TaskUserRelation{
			Tid:      r.TID,
			TaskId:   r.TaskID,
			UserId:   r.UserID,
			Status:   r.Status,
			TaskRole: r.TaskRole,
			TotalNum: r.TotalNum,
		}
		switch r.TaskRole {
		case utils.TASK_PAGE_TAB_ANNOTAT:
			if n, ok := req.NumAnnot[r.UserID]; ok && n != r.TotalNum {
				rel.TotalNum = n
			}
		case utils.TASK_PAGE_TAB_AUDIT:
			if n, ok := req.NumAudit[r.UserID]; ok && n != r.TotalNum {
				rel.TotalNum = n
			}
		case utils.TASK_PAGE_TAB_ANALYS:
			if n, ok := req.NumAnal[r.UserID]; ok && n != r.TotalNum {
				rel.TotalNum = n
			}
		}
		rels = append(rels, rel)
	}
	err = dao.CreateTaskUserRelationInBatches(rels)
	return err
}

func RedoAllocatTaskSingle(req RedoAllocatTaskSingleDTO, mode string) error {
	//ids
	ids, err := dao.GetTcIDsByTaskID(req.TaskId)
	if err != nil {
		return err
	}

	//task calls表更新
	k := 0
	k1 := 0
	for _, a1 := range req.Operator {
		if a1.Num == 0 {
			continue
		}
		k1 += a1.Num
		upd := map[string]interface{}{
			//"status": initStatus,  //不更新状态
		}

		switch req.Type {
		case utils.TASK_PAGE_TAB_ANNOTAT:
			upd["multi_annotator"] = a1.UserID
		case utils.TASK_PAGE_TAB_AUDIT:
			upd["multi_auditor"] = a1.UserID
		case utils.TASK_PAGE_TAB_ANALYS:
			upd["multi_analysts"] = a1.UserID
		}

		_, err := dao.UpdateTaskCallByTcIDs(ids[k:k1], upd)
		if err != nil {
			seelog.Errorf("UpdateTaskCallByTcIDs Annotator data:%+v, err:%s", a1, err.Error())
		}
		k = k1
	}

	tk := dao.TaskInfo{
		Id: req.Id,
	}
	var column string
	var CallIdDoneStartStatus string

	switch req.Type {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		tk.MultiAnnotator = req.TaskOperator
		column = "multi_annotator"
		CallIdDoneStartStatus = utils.TASK_STATUS_ANNOTATED

	case utils.TASK_PAGE_TAB_AUDIT:
		tk.MultiAuditor = req.TaskOperator
		column = "multi_auditor"
		CallIdDoneStartStatus = utils.TASK_STATUS_AUDITED

	case utils.TASK_PAGE_TAB_ANALYS:
		tk.MultiAnalysts = req.TaskOperator
		column = "multi_analysts"
		CallIdDoneStartStatus = utils.TASK_STATUS_ANALYSTED

	}

	_, err = tk.UpdateTaskInfoByID(column)

	//redo 需要删除历史task_user_relation
	err = dao.DeleteTaskUserRelationByRole(req.TaskId, req.Type)
	if err != nil {
		return err
	}

	//任务 用户 角色 关系表 入库
	rels := []dao.TaskUserRelation{}
	for _, r := range req.Relation {
		rel := dao.TaskUserRelation{
			Tid:      r.TID,
			TaskId:   r.TaskID,
			UserId:   r.UserID,
			Status:   r.Status,
			TotalNum: r.TotalNum,
			TaskRole: r.TaskRole,
		}

		if n, ok := req.NumAudit[r.UserID]; ok && n != r.TotalNum {
			rel.TotalNum = n
		}

		//重新计算最新的relation状态, 直接创建对应的doneNum, status
		doneNum, err := dao.GetTaskCallStatusNumByStatus(
			r.TaskID,
			req.Type,
			r.UserID,
			utils.GetCallIdRemainStatus(CallIdDoneStartStatus),
		)
		if err != nil {
			return err
		}
		rel.DoneNum = int(doneNum)
		rel.Status = utils.TASK_STATUS_ALLOCATED
		switch req.Type {
		case utils.TASK_PAGE_TAB_ANNOTAT:
			if int(doneNum) == rel.TotalNum {
				rel.Status = utils.TASK_STATUS_ANNOTATED
			}
		case utils.TASK_PAGE_TAB_AUDIT:
			if int(doneNum) == rel.TotalNum {
				rel.Status = utils.TASK_STATUS_AUDITED
			}
		case utils.TASK_PAGE_TAB_ANALYS:
			if int(doneNum) == rel.TotalNum {
				rel.Status = utils.TASK_STATUS_ANALYSTED
			}
		}

		rels = append(rels, rel)
	}
	err = dao.CreateTaskUserRelationInBatches(rels)
	return err
}

func GetTaskCallsList(taskId string, pageTab string, userID uint32, status []string, pageNum, pageSize int, lastID int32) ([]CallIdStatus, TaskLastLocation, int64, error) {
	var calls = make([]CallIdStatus, 0)
	var total int64
	var loc = TaskLastLocation{}
	taskcalls, total, err := dao.GetTaskCallsByTaskId(taskId, pageTab, userID, status, pageNum, pageSize)
	if err != nil {
		return calls, loc, total, err
	}
	for _, tc := range *taskcalls {
		cs := CallIdStatus{
			CallId:       tc.CallId,
			SerialNumber: tc.SerialNumber,
			Status:       tc.Status,
		}
		calls = append(calls, cs)
		if tc.SerialNumber == int(lastID) {
			loc = TaskLastLocation{
				CallIdStatus: cs,
				Page:         pageNum,
			}
		}
	}
	return calls, loc, total, nil
}

func GetOriginTaskCallsList(
	taskId string,
	pageTab string,
	userID uint32,
	status []string,
	pageNum, pageSize int,
) ([]dao.TaskCall, int64, error) {
	calls, total, err := dao.GetTaskCallsByTaskId(taskId, pageTab, userID, status, pageNum, pageSize)
	if err != nil {
		return *calls, total, err
	}
	return *calls, total, nil
}

// 获取用户在某个任务中最后操作的serial number
func GetTaskUserRelationLastID(taskId string, userID uint32, taskRole string) int32 {
	tur, err := dao.GetTaskUserRelationInfo(taskId, userID, taskRole)
	if err != nil || tur == nil {
		return 0
	}
	return tur.LastId
}

// 获取用户在某个任务中最后操作的serial number 在第几页
func GetTaskUserLastPage(taskID, pageTab string, userID uint32, lastID int32, pageSize int) int {
	total, err := dao.GetTaskCallNumByLastID(taskID, pageTab, userID, lastID)
	if err != nil || total == 0 {
		return 0
	}
	return int(total)/pageSize + 1
}

// 更改task calls 某个callid的状态，newstatus新状态，oldstatus是老状态，为空表示不判断当前状态
func UpdateTaskCallStatus(taskID, callID, newStatus, oldStatus string) (int64, error) {
	return dao.UpdateTaskCallStatus(taskID, callID, newStatus, oldStatus)
}

// 获取任务下callids 某个状态的数量及callids总数量.
func GetTaskStatusNumList(taskID string, status []string, pageTab string, userID uint32) (int, int) {
	var total, num int
	list, err := dao.GetTaskCallStatusNum(taskID, pageTab, userID)
	if err != nil {
		seelog.Errorf("GetTaskStatusNumList err:%s", err.Error())
		return total, num
	}
	for k, v := range list {
		total += v
		if utils.IsInStingSlice(k, status) {
			num += v
		}
	}
	return total, num
}

// 根据条件更新 taskuserrelation表的数据
func UpdateTaskUserRelationByWhere(taskID string, userID uint32, taskRole string, upd map[string]interface{}) error {
	_, err := dao.UpdateTaskUserRelationByWhere(taskID, userID, taskRole, upd)
	return err
}

// 任务预览-calls列表
func GetTaskPreviewCallsList(taskId string, pageTab string, userID uint32, status []string, pageNum, pageSize int) (map[string][]LabelWorkInfoDTO, int64, error) {
	var calls = make(map[string][]LabelWorkInfoDTO, 0)
	var total int64
	taskcalls, total, err := dao.GetTaskCallsByTaskId(taskId, pageTab, userID, status, pageNum, pageSize)
	if err != nil {
		return calls, total, err
	}
	for _, tc := range *taskcalls {
		//获取callid 对应的labels
		labels := GetOneCallLabelWork(taskId, tc.CallId, pageTab)
		calls[tc.CallId] = labels
	}
	return calls, total, nil
}

// 生成任务Excel文件并返回下载地址
func GetTaskDownloadFileAddr(taskId string) (string, error) {
	seelog.Infof("taskID:%s, startdownload", taskId)

	calls, _, err := dao.GetTaskCallsByTaskId(taskId, "", 0, nil, 0, 0)
	seelog.Infof("taskID:%s, GetTaskCallsByTaskId", taskId)

	if err != nil {
		seelog.Errorf("get %s distribute faild:%s", taskId, err)
		return "", err
	}
	if len(*calls) == 0 {
		seelog.Errorf("no content for task:%s", taskId)
		return "", errors.New("no content for task:" + taskId)
	}

	var excelTable = make([]TaskCallLabelListItem, 0, len(*calls))
	headerList := []string{
		"call_id",
	}
	headerKeys := false

	for _, tc := range *calls {
		labels := GetOneCallLabelWork(taskId, tc.CallId, "")
		if len(labels) == 0 {
			continue
		}
		if !headerKeys {
			for _, lab := range labels {
				headerList = append(headerList, lab.LabelName)
			}
			headerKeys = true
		}
		item := TaskCallLabelListItem{
			CallID: tc.CallId,
			Labels: labels,
		}
		excelTable = append(excelTable, item)

	}
	seelog.Infof("taskID:%s, excelTable", taskId)
	// filename := fileprocess.CFG_LOCALREPORTFILEPATH + "/" + taskId + ".xlsx"
	// // xlserr := exportTaskFile(excelTable, headerList, filename)
	// if xlserr != nil {
	// 	seelog.Errorf("generate xls file failed %s", xlserr.Error())
	// 	return "", xlserr
	// }
	seelog.Infof("taskID:%s, exportTaskFile", taskId)

	return fileprocess.CFG_DOWNLOADURL + taskId + ".xlsx", nil
}

// 生成任务Excel文件并返回下载地址
func GetTaskDownloadFileAddrV2(taskId string) (string, error) {
	seelog.Infof("taskID:%s, startdownload", taskId)

	calls, err := dao.GetCallIDsByTaskID(taskId)
	seelog.Infof("taskID:%s, GetTaskCallsByTaskId", taskId)

	if err != nil {
		seelog.Errorf("get %s distribute faild:%s", taskId, err)
		return "", err
	}
	clen := len(calls)
	if clen == 0 {
		seelog.Errorf("no content for task:%s", taskId)
		return "", errors.New("no content for task:" + taskId)
	}
	excelTitleArr := []string{
		"RobotType",
		"RobotName",
		"StartTime",
		"CallID",
		"BillSec",
		"Intention",
		"TalkRound",
		"Sentence1",
		"Final Intention",
		"Check",
		"Problem",
		"Remark",
		"Chinese",
		"Dialogue Fluency",
		"Real Intention",
		"Speech Recognition",
		"Semantic Comprehension",
		"Robot language",
		"Robot Reaction Speed",
		"Noise situation",
		"Robot Recognition",
		"User cooperation",
		"Problem2",
		"Remark2",
	}
	var excelTable = make([]TaskCallLabelListDown, 0, clen)
	headerList := []string{}
	headerListLost := []string{}
	labels := GetOneCallLabelWork(taskId, calls[0], "")
	var labNames = map[string]string{}
	if len(labels) == 0 {
		seelog.Errorf("no label for task:%s", taskId)
		return "", errors.New("no label for task:" + taskId)
	}
	for _, lab := range labels {
		labNames[lab.LabelName] = lab.LabelName
		headerListLost = append(headerListLost, lab.LabelName)
	}
	for _, et := range excelTitleArr {
		if na, ok := labNames[et]; ok {
			headerList = append(headerList, na)
			delete(labNames, na)
		}
	}
	for _, et := range headerListLost {
		if et == "Sentence" {
			continue
		}
		if _, ok := labNames[et]; ok {
			headerList = append(headerList, et)
		}
	}

	j := 500
	s := 0
	e := 0
	for e < clen {
		e = s + j
		if e > clen {
			e = clen
		}
		cids := calls[s:e]
		list, err := dao.GetLabelWorkDetailByCallIDs(taskId, cids)
		if err != nil {
			seelog.Errorf("get %s GetLabelWorkDetailByCallIDs faild:%s", taskId, err.Error())
			continue
		}

		for _, cid := range cids {
			cinfo, ok := list[cid]
			if !ok {
				continue
			}
			labels := make([]LabelWorkInfoDown, 0)
			for _, name := range headerList {
				if nv, ok := cinfo[name]; ok {
					isColor := 0
					if labInfo, ok := LabelInfoCache[nv.LabelId]; ok {
						isColor = labInfo.IsColor
					}

					nd := LabelWorkInfoDown{
						LabelId:        nv.LabelId,
						LabelName:      nv.LabelName,
						IsColor:        isColor,
						AuditorContent: nv.AuditorContent,
						LabelValue:     nv.LabelValue,
					}
					labels = append(labels, nd)
				}
			}
			item := TaskCallLabelListDown{
				CallID: cid,
				Labels: labels,
			}
			excelTable = append(excelTable, item)
		}
		s += j
	}

	seelog.Infof("taskID:%s, excelTable", taskId)
	filename := fileprocess.CFG_LOCALREPORTFILEPATH + "/" + taskId + ".xlsx"
	xlserr := exportTaskFile(excelTable, headerList, filename)
	if xlserr != nil {
		seelog.Errorf("generate xls file failed %s", xlserr.Error())
		return "", xlserr
	}
	seelog.Infof("taskID:%s, exportTaskFile", taskId)

	return fileprocess.CFG_DOWNLOADURL + taskId + ".xlsx", nil
}

// 根据taskid，role，status, 获取任务下用户所属内容完成情况，用于判断任务是否进入下一阶段
func GetTaskCurrentStatusIsDone(taskID, taskRole, status string) (bool, error) {
	list, err := dao.GetTaskUserRelationStatusNum(taskID, taskRole)
	if err != nil {
		return false, err
	}
	var total int
	var snum int
	for st, n := range list {
		if st == status {
			snum = n
		}
		total += n
	}
	return total == snum, nil
}

// 根据taskid，role, 获取任务下用户所属内容完成数量
func GetTaskUserDoneNumByTab(taskID, taskRole string) (int, int, error) {
	list, err := dao.GetTaskUserRelationDoneNum(taskID, taskRole)
	if err != nil {
		return 0, 0, err
	}
	var total int
	var snum int
	for _, n := range *list {
		snum += n.DoneNum
		total += n.TotalNum
	}
	return total, snum, nil
}

//
//func ChangeTaskNextStatus(taskID string) (string, error) {
//
//	for {
//		nextStatus, err := GetTaskNextStatus(taskID)
//		if err != nil {
//			return "", err
//		}
//		//计算需要切换的status， 是否已经完成
//
//		switch nextStatus {
//
//		case utils.TASK_STATUS_ANNOTATED:
//			err = UpdateTaskStatusByTaskID(taskID, nextStatus)
//			if err != nil {
//				return "", err
//			}
//
//			//已流转到待审核： 验证是否存在待审核的callid，  不存在需要继续流转
//			TaskLeftCallids, err := GetUserTaskCallsByStatus(
//				taskID,
//				0,
//				"",
//				[]string{
//					utils.TASK_STATUS_ALLOCATED,
//					utils.TASK_STATUS_ANNOTATED,
//					utils.CALLID_STATUS_PRE_AUDIT,
//				},
//			)
//			if err != nil {
//				return "", err
//			}
//
//		case utils.TASK_STATUS_AUDITED:
//			//待分析： 验证是否存在待分析的callid
//
//		case utils.TASK_STATUS_COMPLETED:
//			//已完成：直接更新db， return
//
//		}
//
//	}
//
//}

func GetCallIDNextStatus(taskinfo *TaskInfo, currentStatus string) (string, error) {

	if taskinfo == nil {
		return "", errors.New("taskinfo is nil")
	}
	var status string
	switch currentStatus {
	case utils.TASK_STATUS_ANNOTATED:
		status = utils.CALLID_STATUS_PRE_AUDIT
		if len(taskinfo.Auditor) > 0 {
			break
		}
		if len(taskinfo.Analysts) > 0 {
			status = utils.TASK_STATUS_PRE_ANALYST
			break
		}
		status = utils.TASK_STATUS_COMPLETED

	case utils.TASK_STATUS_AUDITED:
		status = utils.TASK_STATUS_PRE_ANALYST
		if len(taskinfo.Analysts) > 0 {
			break
		}
		status = utils.TASK_STATUS_COMPLETED
	}
	return status, nil
}

// 获取任务下一个状态.
func GetTaskNextStatus(taskID string) (string, error) {
	taskinfo, err := dao.GetTaskInfoByID(taskID)
	var status string
	if err != nil {
		seelog.Errorf("GetTaskInfoByID failed:%s", err.Error())
		return status, err
	}
	if taskinfo == nil {
		return status, nil
	}
	switch taskinfo.Status {
	case utils.TASK_STATUS_ANNOTATING, utils.TASK_STATUS_ALLOCATED:
		status = utils.TASK_STATUS_ANNOTATED
		if taskinfo.MultiAuditor != "" {
			status = utils.TASK_STATUS_ANNOTATED
			break
		}
		if taskinfo.MultiAnalysts != "" {
			status = utils.TASK_STATUS_AUDITED
			break
		}
		status = utils.TASK_STATUS_COMPLETED

	case utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_AUDITING:
		status = utils.TASK_STATUS_AUDITED
		if taskinfo.MultiAnalysts != "" {
			status = utils.TASK_STATUS_AUDITED
			break
		}
		status = utils.TASK_STATUS_COMPLETED

	case utils.TASK_STATUS_AUDITED, utils.TASK_STATUS_ANALYSTING:
		status = utils.TASK_STATUS_COMPLETED
	}
	return status, nil
}

// 根据taskid ,old status, 更新task status
func UpdateTaskStatusByTaskStatus(taskid, newStatus, oldStatus string) error {
	return dao.UpdateTaskStatusByTaskStatus(taskid, newStatus, oldStatus)
}

// 根据更新label work信息，更新task user relation数据
func UpdateTaskUserRelationForLabel(taskid, callid, taskRole string, userID uint32, serialNumber int) error {
	relInfo, err := dao.GetTaskUserRelationInfo(taskid, userID, taskRole)
	if err != nil {
		return err
	}
	if relInfo == nil {
		return nil
	}
	relInfo.LastId = int32(serialNumber)
	columns := []string{"last_id"}
	switch taskRole {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		if relInfo.Status == utils.TASK_STATUS_ALLOCATED {
			columns = append(columns, "status")
			relInfo.Status = utils.TASK_STATUS_ANNOTATING
		}
	case utils.TASK_PAGE_TAB_AUDIT:
		if relInfo.Status == utils.TASK_STATUS_ANNOTATED {
			columns = append(columns, "status")
			relInfo.Status = utils.TASK_STATUS_AUDITING
		}
	case utils.TASK_PAGE_TAB_ANALYS:
		if relInfo.Status == utils.TASK_STATUS_AUDITED {
			columns = append(columns, "status")
			relInfo.Status = utils.TASK_STATUS_ANALYSTING
		}
	}
	_, err = relInfo.UpdateTaskUserRelationByID(columns)
	return err
}

// 根据条件更新 taskcall表的数据
func UpdateTaskCallByWhere(taskID, callID string, upd map[string]interface{}) error {
	_, err := dao.UpdateTaskCallByWhere(taskID, callID, upd)
	return err
}

// 根据taskid callid 获取详情
func GetTaskCallInfoByTaskIDCallID(taskID, callID string) (TaskCallItem, error) {
	var tc = TaskCallItem{}
	result, err := dao.GetTaskCallInfoByTaskIDCallID(taskID, callID)
	if err != nil || result == nil {
		return tc, err
	}
	mannotator, _ := strconv.ParseInt(result.MultiAnnotator, 10, 64)
	mauditor, _ := strconv.ParseInt(result.MultiAuditor, 10, 64)
	manalysts, _ := strconv.ParseInt(result.MultiAnalysts, 10, 64)
	annotator, _ := strconv.ParseInt(result.ActualAnnotator, 10, 64)
	auditor, _ := strconv.ParseInt(result.ActualAuditor, 10, 64)
	analysts, _ := strconv.ParseInt(result.ActualAnalysts, 10, 64)
	tc = TaskCallItem{
		TcId:            result.TcId,
		TaskId:          result.TaskId,
		CallId:          result.CallId,
		MultiAnnotator:  uint32(mannotator),
		MultiAuditor:    uint32(mauditor),
		MultiAnalysts:   uint32(manalysts),
		ActualAnnotator: uint32(annotator),
		ActualAuditor:   uint32(auditor),
		ActualAnalysts:  uint32(analysts),
		SerialNumber:    result.SerialNumber,
		Status:          result.Status,
		RejectReason:    result.RejectReason,
	}
	return tc, nil
}

// 根据taskuserrelation 获取详情
func GetTaskUserRelationInfo(taskID, taskRole string, userID uint32) (TaskUserRelationItem, error) {
	var tc = TaskUserRelationItem{}
	result, err := dao.GetTaskUserRelationInfo(taskID, userID, taskRole)
	if err != nil || result == nil {
		return tc, err
	}
	tc = TaskUserRelationItem{
		TurID:    result.TurId,
		TID:      result.TurId,
		TaskID:   result.TaskId,
		UserID:   result.UserId,
		TaskRole: result.TaskRole,
		TotalNum: result.TotalNum,
		Status:   result.Status,
	}
	return tc, nil
}

// 清除任务本地音频文件
func ClearTaskLocalWavFile(taskID string) {
	subPath := path.Join(fileprocess.CFG_LOCALRECORDINGFILEPATH, taskID)
	if err := os.RemoveAll(subPath); err != nil {
		seelog.Errorf("ClearTaskLocalWavFile err:%s, path:%s", err.Error(), subPath)
	}
}

func GetTaskProgressDetails(ctx context.Context, taskID string) ([]*dao.TaskProgressDetail, error) {
	details, err := dao.GetTaskProgressDetails(taskID)
	if err != nil {
		return details, err
	}

	detailsMap := make(map[string]*dao.TaskProgressDetail)

	for _, detail := range details {
		detailsMap[detail.Status] = detail
	}

	//追加任务状态num=0的数据
	var newDetails []*dao.TaskProgressDetail
	for _, s := range utils.CallIdStatusFlow {
		if detailsMap[s] == nil {
			newDetails = append(newDetails, &dao.TaskProgressDetail{
				Status: s,
				Num:    0,
			})
		} else {
			newDetails = append(newDetails, detailsMap[s])
		}
	}

	return newDetails, nil
}
