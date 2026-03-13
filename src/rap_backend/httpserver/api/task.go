package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rap_backend/fileprocess"
	"rap_backend/internal/context"
	"rap_backend/service"
	"rap_backend/utils"
	"reflect"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func TaskDetail(ctx *gin.Context) {
	input := service.TaskDetailDTO{}
	if err := ctx.Bind(&input); err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task list request marsh error", nil))
		return
	}
	userinfo := context.GetUserInfo(ctx)
	if userinfo == nil {
		seelog.Errorf("get task list userinfo is null ")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "get task list userinfo is null", nil))
		return
	}

	info, err := service.GetOriginTaskInfoByID(input.TaskId, true)
	if err != nil {
		seelog.Errorf("get task info error failed:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get task info error failed", nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(info))
}

func SubmitStatistics(ctx *gin.Context) {
	input := service.SubmitStatisticsDTO{}
	if err := ctx.Bind(&input); err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task list request marsh error", nil))
		return
	}
	userinfo := context.GetUserInfo(ctx)
	if userinfo == nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "get task list userinfo is null", nil))
		return
	}

	var status []string
	switch input.Type {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		status = []string{utils.TASK_STATUS_ANNOTATED}
	case utils.TASK_PAGE_TAB_AUDIT:
		status = []string{utils.TASK_STATUS_AUDITED}
	case utils.TASK_PAGE_TAB_ANALYS:
		status = []string{utils.TASK_STATUS_ANALYSTED}

	default:
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task list request marsh error", nil))
		return
	}

	total, num := service.GetTaskStatusNumList(
		input.TaskId,
		status,
		input.Type,
		userinfo.UserID,
	)

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(gin.H{
		"total": total,
		"done":  num,
	}))
}

// 任务管理页面-任务列表.
func GetTaskList(ctx *gin.Context) {
	seelog.Infof("get task list cmd")
	input := service.GetTaskListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get task list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task list request marsh error", nil))
		return
	}
	userinfo := context.GetUserInfo(ctx)
	if userinfo == nil {
		seelog.Errorf("get task list userinfo is null ")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "get task list userinfo is null", nil))
		return
	}
	seelog.Infof("get task list cmd page num:%d, page size:%d", input.PageNum, input.PageSize)
	req := service.TaskListFilter{
		PageTab:     utils.TASK_PAGE_TAB_MANAGE,
		UserIDs:     []uint32{userinfo.UserID},
		SearchType:  input.SearchType,
		SearchValue: input.SearchValue,
		PageNum:     input.PageNum,
		PageSize:    input.PageSize,
		CurrentUID:  userinfo.UserID,
	}
	if userinfo.DataAll { //全局数据，查看所有
		req.UserIDs = []uint32{}
	}
	if input.SearchValue == "" {
		input.SearchType = ""
	}
	if input.SearchType == utils.TASK_SEARCH_TYPE_OPERATOR {
		uu, err1 := service.GetUserInfoByUserName(input.SearchValue)
		if err1 != nil || uu.UserID == 0 {
			seelog.Errorf("get task list error failed")
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get task list user not exist", nil))
			return
		}
		req.SearchValue = fmt.Sprintf("%d", uu.UserID)
	}
	taskInfoList, cnt, err2 := service.GetTaskListByFilter(req)
	if err2 != nil {
		seelog.Errorf("get task list error failed:%s", err2.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get task list error failed", nil))
		return
	}
	var ret service.GetTaskListRetDTO
	if taskInfoList != nil {
		ret.TaskInfoList = *taskInfoList
	}
	ret.TaskCnt = cnt
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

// 任务分配页面-任务列表.
func GetAllocatTaskList(ctx *gin.Context) {
	seelog.Infof("get GetAllocatTaskList cmd")
	input := service.GetTaskListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get task list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get GetAllocatTaskList request marsh error", nil))
		return
	}
	userinfo := context.GetUserInfo(ctx)
	if userinfo == nil {
		seelog.Errorf("get task list userinfo is null ")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "get GetAllocatTaskList userinfo is null", nil))
		return
	}
	seelog.Infof("get GetAllocatTaskList cmd page num:%d, page size:%d", input.PageNum, input.PageSize)
	req := service.TaskListFilter{
		PageTab:     utils.TASK_PAGE_TAB_ALLOCAT,
		UserIDs:     []uint32{userinfo.UserID},
		SearchType:  input.SearchType,
		SearchValue: input.SearchValue,
		PageNum:     input.PageNum,
		PageSize:    input.PageSize,
	}
	if userinfo.DataAll { //全局数据，查看所有
		req.UserIDs = []uint32{}
	}
	if input.SearchValue == "" {
		input.SearchType = ""
	}
	if input.SearchType == utils.TASK_SEARCH_TYPE_OPERATOR {
		uu, err1 := service.GetUserInfoByUserName(input.SearchValue)
		if err1 != nil || uu.UserID == 0 {
			seelog.Errorf("get GetAllocatTaskList error failed")
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get task list user not exist", nil))
			return
		}
		req.SearchValue = fmt.Sprintf("%d", uu.UserID)
	}
	taskInfoList, cnt, err2 := service.GetTaskListByFilter(req)
	if err2 != nil {
		seelog.Errorf("get GetAllocatTaskList error failed %s", err2.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get GetAllocatTaskList error failed", nil))
		return
	}
	var ret service.GetTaskListRetDTO
	if taskInfoList != nil {
		ret.TaskInfoList = *taskInfoList
	}
	ret.TaskCnt = cnt
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

// 默认选中的labelid
func LabelDefault(ctx *gin.Context) {
	input := service.CheckCallIdDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("check callid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "check callid request marsh error", nil))
		return
	}
	labelNames := []string{
		"RobotType",
		"RobotName",
		"StartTime",
		"CallID",
		"BillSec",
		"Intention",
		"TalkRound",
		"Sentence",
	}

	var ret = service.LabelDefaultRetDTO{
		Labels: make([]string, 0, len(labelNames)),
	}
	for _, name := range labelNames {
		if lab, ok := service.LabelNameInfoCache[name]; ok {
			ret.Labels = append(ret.Labels, lab.LabelId)
		}
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

// 检测callid在gauss callinfo是否存在
func CheckUploadCallId(ctx *gin.Context) {
	seelog.Infof("receive check callid cmd")
	input := service.CheckCallIdDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("check callid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "check callid request marsh error", nil))
		return
	}
	if len(input.Callid) == 0 {
		seelog.Errorf("scheck callid request callid is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid list is null", nil))
		return
	}
	if len(input.Callid) > 5000 {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid list more than 5000", nil))
		return
	}

	seelog.Infof("receive check callid cmd, num:%d", len(input.Callid))
	oklst, exceptlst, err2 := service.CheckUploadCallId(input.Callid)
	if err2 != nil {
		seelog.Errorf("scheck callid failed:%s", err2.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "check callid request failed", nil))
	}
	var ret service.CheckCallIdRetDTO
	ret.CheckedCallid = removeRepeatElement(oklst)
	ret.ExceptionCallid = exceptlst

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

func removeRepeatElement(list []string) []string {
	// 创建一个临时map用来存储数组元素
	temp := make(map[string]struct{})
	index := 0
	// 将元素放入map中
	for _, v := range list {
		temp[v] = struct{}{}
	}
	tempList := make([]string, len(temp))
	for key := range temp {
		tempList[index] = key
		index++
	}
	return tempList
}

//创建任务
// func CreateTask(ctx *gin.Context) {
// 	seelog.Infof("receive create task cmd")
// 	input := service.CreateTaskDTO{}
// 	if err := ctx.Bind(&input); err != nil {
// 		seelog.Errorf("create task request marsh error :%s", err.Error())
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "create task request marsh error", nil))
// 		return
// 	}
// 	if len(input.Callid) == 0 {
// 		seelog.Errorf("create task request callid is null")
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "callid list is null", nil))
// 		return
// 	}

// 	if input.TaskName == "" {
// 		seelog.Errorf("create task request taskname is null")
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "taskname is null", nil))
// 		return
// 	}
// 	err, taskid := service.CreateNewTask(input)
// 	if err != nil {
// 		seelog.Errorf("create task failed:%s", err.Error())
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "create task failed", nil))
// 		return
// 	}
// 	subNum := len(input.Callid) / len(input.TaskOperators)
// 	seelog.Infof("sub num:%d", subNum)
// 	offset := 0
// 	for idx, optor := range input.TaskOperators {
// 		var distributeCallNum = subNum
// 		if idx+1 >= len(input.TaskOperators) && idx != 0 {
// 			distributeCallNum = len(input.Callid) - (idx+1)*subNum
// 		}
// 		seelog.Infof("optor:%s, offset :%d, distributeCallNum:%d", optor, offset, distributeCallNum)
// 		err2, subtaskid := service.CreateNewSubTask(input, taskid, optor, offset, distributeCallNum)
// 		if err2 != nil {
// 			seelog.Errorf("create sub task failed:%s", err2.Error())
// 			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "create sub task failed", nil))
// 			return
// 		}
// 		err3 := service.CreateNewDistribute(taskid, subtaskid, input.TaskCreator)
// 		if err3 != nil {
// 			seelog.Errorf("create distribution failed:%s", err3.Error())
// 			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "create distribution failed", nil))
// 			return
// 		}
// 		go fileprocess.PrepareTaskRecordingFiles(taskid, subtaskid, input.Callid)
// 	}
// 	var ret service.CreateTaskRetDTO
// 	ret.Taskid = taskid

// 	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
// }

// 创建任务
func CreateTaskV2(ctx *gin.Context) {
	seelog.Infof("receive create task cmd")
	input := service.CreateTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("create task request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "create task request marsh error", nil))
		return
	}
	if len(input.Callid) == 0 || input.TaskName == "" || len(input.TaskOperators) == 0 || input.FinishTime == "" || len(input.LabelId) == 0 {
		seelog.Errorf("create task request param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	//检查labelid是否都正确
	for _, labID := range input.LabelId {
		if _, ok := service.LabelInfoCache[labID]; !ok {
			seelog.Errorf("labelid not exist %s", labID)
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "labelid not exist", nil))
			return
		}
	}

	task, err := service.CheckTaskName(input.TaskName)
	if err != nil || task != "" {
		seelog.Errorf("taskname exist %s", input.TaskName)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_TASK_NAME_EXIST, "", nil))
		return
	}
	userinfo := context.GetUserInfo(ctx)
	if userinfo == nil {
		seelog.Errorf("create task userinfo is nil")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "create task userinfo is nil", nil))
		return
	}
	country := service.GetCountryInfoByUID(userinfo.UserID)
	finishTime := utils.TimeZone2UTCTime(input.FinishTime, country.ZoneId)
	if finishTime.Before(utils.NowUTC()) {
		seelog.Errorf("create task finishTime less than current time")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "finishTime less than current time", nil))
		return
	}
	input.StartTimer = utils.NowUTC()
	input.FinishTimer = finishTime

	input.CreatorId = userinfo.UserID
	input.TaskCreator = userinfo.UserName
	err, taskid := service.CreateNewTask(input)
	if err != nil {
		seelog.Errorf("create task failed:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task failed", nil))
		return
	}

	err = service.CreateTaskCalls(taskid, input.Callid)
	if err != nil {
		seelog.Errorf("create task calls failed:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task calls failed", nil))
		return
	}
	err = service.CreateTaskLabel2(input, taskid)
	go fileprocess.PrepareTaskRecordingFiles(taskid, input.Callid)

	var ret service.CreateTaskRetDTO
	ret.Taskid = taskid

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

func GetTaskDetailByTaskId(ctx *gin.Context) {
	seelog.Infof("receive get task  detail by taskid cmd")
	input := service.GetTaskCallIdListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get task detail by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task detail by taskid request marsh error", nil))
		return
	}
	if input.TaskId == "" {
		seelog.Errorf("get task detail by taskid, taskid is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task id is null", nil))
	}

	if input.PageNum == 0 {
		seelog.Warnf("get task detail request page num is 0")
		input.PageNum = 1
	}
	if input.PageSize == 0 {
		seelog.Warnf("get task detail request page size is 0")
		input.PageSize = 20
	}

	getTaskDetailRet, err := service.GetTaskCallIdList(input.TaskId, input.SubTaskId, input.PageNum, input.PageSize)
	if err != nil {
		seelog.Errorf("get task detail by taskid failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task failed", nil))
		return
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(getTaskDetailRet))

}

// func GetSubTaskStatics(ctx *gin.Context) {
// 	seelog.Infof("receive get sub task statics by subtaskid cmd")
// 	input := service.GetSubTaskStaticsDTO{}
// 	if err := ctx.Bind(&input); err != nil {
// 		seelog.Errorf("get sub task statics by subtaskid request marsh error :%s", err.Error())
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "get sub task statics by subtaskid request marsh error", nil))
// 		return
// 	}
// 	if input.SubTaskId == "" {
// 		seelog.Errorf("get sub task statics by subtaskid, subtaskid is null")
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "sub task id is null", nil))
// 	}
// 	ret := service.GetSubTaskStatics(input.SubTaskId)
// 	if ret == nil {
// 		ctx.JSON(http.StatusOK, NewCommonSuccessResp(nil))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
// }

func DownloadTask(ctx *gin.Context) {
	seelog.Infof("receive get task  detail by taskid cmd")
	input := service.DownloadTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("DownloadTask by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "DownloadTask by taskid request marsh error", nil))
		return
	}
	if input.TaskId == "" {
		seelog.Errorf("DownloadTask by taskid, taskid and subtask id is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid and subtask id is null", nil))
	}
	var err2 error
	var downloadUrl string
	downloadUrl, err2 = service.GetTaskDownloadFileAddrV2(input.TaskId)
	if err2 != nil {
		seelog.Errorf("create %s download file failed :%s", input.TaskId, err2.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "download task file failed", nil))
		return
	}
	var resp = service.DownloadTaskRetDTO{
		DownloadUrl: downloadUrl,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func DeleteTask(ctx *gin.Context) {
	seelog.Infof("receive get task  detail by taskid cmd")
	input := service.DownloadTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("DeleteTask by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "DeleteTask by taskid request marsh error", nil))
		return
	}
	if input.TaskId == "" {
		seelog.Errorf("DeleteTask by taskid, taskid and subtask id is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid is null", nil))
		return
	}

	err3 := service.UpdateTaskStatusByTaskID(input.TaskId, utils.TASK_STATUS_DELETED)
	if err3 != nil {
		seelog.Errorf("DeleteTask %s failed :%s", input.TaskId, err3.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "delete task failed", nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func ReDoAllocatTask(ctx *gin.Context) {
	input := service.DoAllocatTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("RedoAllocatTask by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "DeleteTask by taskid request marsh error", nil))
		return
	}
	logData, _ := json.Marshal(input)
	seelog.Infof("receive ReDoAllocatTask: %v", logData)

	userID := context.GetUID(ctx)
	if input.TaskId == "" || (len(input.Annotator) == 0 && len(input.Auditor) == 0 && len(input.Analysts) == 0) {
		seelog.Errorf("RedoAllocatTask by taskid, param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("RedoAllocatTask by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "taskid is error", nil))
		return
	}

	if taskinfo.TaskStatus == utils.TASK_STATUS_CREATED || taskinfo.TaskStatus == utils.TASK_STATUS_COMPLETED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid status error", nil))
		return
	}
	if !utils.IsInUint32Slice(taskinfo.Allocator, userID) {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot assign this task", nil))
		return
	}
	totalAnnot := 0
	userAnnot := []uint32{}
	numAnnot := make(map[uint32]int, 0)
	for _, m := range input.Annotator {
		if m.Num == 0 {
			continue
		}
		totalAnnot += m.Num
		if _, nok := numAnnot[m.UserID]; nok {
			numAnnot[m.UserID] += m.Num
			continue
		}
		numAnnot[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_ANNOTAT,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAnnot = append(userAnnot, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.StatusInit = utils.TASK_STATUS_ALLOCATED
		}
	}
	input.NumAnnot = numAnnot
	totalAuditor := 0
	userAuditor := []uint32{}
	numAudit := make(map[uint32]int, 0)
	for _, m := range input.Auditor {
		if m.Num == 0 {
			continue
		}
		totalAuditor += m.Num
		if _, nok := numAudit[m.UserID]; nok {
			numAudit[m.UserID] += m.Num
			continue
		}
		numAudit[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_AUDIT,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAuditor = append(userAuditor, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.StatusInit = utils.CALLID_STATUS_PRE_AUDIT
			input.TaskStatusInit = utils.TASK_STATUS_ANNOTATED
		}
	}
	input.NumAudit = numAudit
	totalAnalysts := 0
	userAnalysts := []uint32{}
	numAnal := make(map[uint32]int, 0)
	for _, m := range input.Analysts {
		if m.Num == 0 {
			continue
		}
		totalAnalysts += m.Num
		if _, nok := numAnal[m.UserID]; nok {
			numAnal[m.UserID] += m.Num
			continue
		}
		numAnal[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_ANALYS,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAnalysts = append(userAnalysts, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.TaskStatusInit = utils.TASK_STATUS_AUDITED
			input.StatusInit = utils.TASK_STATUS_PRE_ANALYST
		}
	}
	input.NumAnal = numAnal
	if input.StatusInit == "" {
		seelog.Errorf("RedoAllocatTask by taskid, 标注人，审核人，结果分析人至少选一个")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "Select at least one of the annotator, the reviewer, and the result analyst", nil))
		return
	}
	if (totalAnalysts > 0 && totalAnalysts != taskinfo.CallNum) || (totalAuditor > 0 && totalAuditor != taskinfo.CallNum) || (totalAnnot > 0 && totalAnnot != taskinfo.CallNum) {
		seelog.Errorf("RedoAllocatTask by taskid, CallNum is no match")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "CallNum is no match", nil))
		return
	}
	//分配任务
	input.Id = taskinfo.Id
	input.TaskAnnotator = utils.UInt32SliceToString(userAnnot)
	input.TaskAuditor = utils.UInt32SliceToString(userAuditor)
	input.TaskAnalysts = utils.UInt32SliceToString(userAnalysts)
	err3 := service.DoAllocatTask(input, "redo")
	if err3 != nil {
		seelog.Errorf("ReDoAllocatTask %s failed :%s", input.TaskId, err3.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "DoAllocatTask task failed", nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func RedoAllocatTaskSingle(ctx *gin.Context) {
	input := service.RedoAllocatTaskSingleDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("RedoAllocatTask by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "DeleteTask by taskid request marsh error", nil))
		return
	}
	logData, _ := json.Marshal(input)
	seelog.Infof("receive ReDoAllocatTask: %v", string(logData))

	userID := context.GetUID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "not login", nil))
	}

	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("RedoAllocatTask by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "taskid is error", nil))
		return
	}

	if taskinfo.TaskStatus == utils.TASK_STATUS_CREATED || taskinfo.TaskStatus == utils.TASK_STATUS_COMPLETED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid status error", nil))
		return
	}
	if !utils.IsInUint32Slice(taskinfo.Allocator, userID) {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot assign this task", nil))
		return
	}

	//task已经走过的状态，不再允许重新分配
	// 获取另外两个角色的数据，不能同时为空
	switch input.Type {
	case utils.TASK_PAGE_TAB_ANALYS:
		if len(taskinfo.Annotator) == 0 && len(taskinfo.Auditor) == 0 {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "At least one process exists", nil))
			return
		}

		if reflect.DeepEqual(taskinfo.Analysts, input.Operator) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "no adjustment required", nil))
			return
		}
		if !utils.IsAfterCurrentTaskStatus(taskinfo.TaskStatus, utils.TASK_STATUS_ANALYSTING) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "completed processes cannot be reassigned", nil))
			return
		}

	case utils.TASK_PAGE_TAB_AUDIT:
		if len(taskinfo.Annotator) == 0 && len(taskinfo.Analysts) == 0 {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "At least one process exists", nil))
			return
		}
		if reflect.DeepEqual(taskinfo.Auditor, input.Operator) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "no adjustment required", nil))
			return
		}
		if !utils.IsAfterCurrentTaskStatus(taskinfo.TaskStatus, utils.TASK_STATUS_AUDITING) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "completed processes cannot be reassigned", nil))
			return
		}
	case utils.TASK_PAGE_TAB_ANNOTAT:
		if len(taskinfo.Auditor) == 0 && len(taskinfo.Analysts) == 0 {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "At least one process exists", nil))
			return
		}
		if reflect.DeepEqual(taskinfo.Annotator, input.Operator) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "no adjustment required", nil))
			return
		}
		if !utils.IsAfterCurrentTaskStatus(taskinfo.TaskStatus, utils.TASK_STATUS_ANNOTATING) {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "completed processes cannot be reassigned", nil))
			return
		}
	}

	totalAuditor := 0
	userAuditor := []uint32{}
	numAudit := make(map[uint32]int, 0)
	for _, m := range input.Operator {
		if m.Num == 0 {
			continue
		}
		totalAuditor += m.Num
		if _, nok := numAudit[m.UserID]; nok {
			numAudit[m.UserID] += m.Num
			continue
		}
		numAudit[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: input.Type,
			TotalNum: m.Num,
			//Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAuditor = append(userAuditor, m.UserID)
		input.Relation = append(input.Relation, rel)
	}
	input.NumAudit = numAudit

	if totalAuditor > 0 && totalAuditor != taskinfo.CallNum {
		seelog.Errorf("RedoAllocatTask by taskid, CallNum is no match")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "CallNum is no match", nil))
		return
	}
	//分配任务
	input.Id = taskinfo.Id
	input.TaskOperator = utils.UInt32SliceToString(userAuditor)
	if err := service.RedoAllocatTaskSingle(input, "redo"); err != nil {
		seelog.Errorf("ReDoAllocatTask %s failed :%s", input.TaskId, err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "DoAllocatTask task failed", nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 任务分配
func DoAllocatTask(ctx *gin.Context) {
	seelog.Infof("receive get task  detail by taskid cmd")
	input := service.DoAllocatTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("DoAllocatTask by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "DeleteTask by taskid request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)
	if input.TaskId == "" || (len(input.Annotator) == 0 && len(input.Auditor) == 0 && len(input.Analysts) == 0) {
		seelog.Errorf("DoAllocatTask by taskid, param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("DoAllocatTask by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "taskid is error", nil))
		return
	}
	if taskinfo.TaskStatus == utils.TASK_STATUS_ALLOCATED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task has been assigned", nil))
		return
	}
	if taskinfo.TaskStatus != utils.TASK_STATUS_CREATED {
		seelog.Errorf("DoAllocatTask by taskid, taskid status error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid status error", nil))
		return
	}
	if !utils.IsInUint32Slice(taskinfo.Allocator, userID) {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot assign this task", nil))
		return
	}
	totalAnnot := 0
	userAnnot := []uint32{}
	numAnnot := make(map[uint32]int, 0)
	for _, m := range input.Annotator {
		if m.Num == 0 {
			continue
		}
		totalAnnot += m.Num
		if _, nok := numAnnot[m.UserID]; nok {
			numAnnot[m.UserID] += m.Num
			continue
		}
		numAnnot[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_ANNOTAT,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAnnot = append(userAnnot, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.StatusInit = utils.TASK_STATUS_ALLOCATED
			input.TaskStatusInit = utils.TASK_STATUS_ALLOCATED
		}
	}
	input.NumAnnot = numAnnot
	totalAuditor := 0
	userAuditor := []uint32{}
	numAudit := make(map[uint32]int, 0)
	for _, m := range input.Auditor {
		if m.Num == 0 {
			continue
		}
		totalAuditor += m.Num
		if _, nok := numAudit[m.UserID]; nok {
			numAudit[m.UserID] += m.Num
			continue
		}
		numAudit[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_AUDIT,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAuditor = append(userAuditor, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.StatusInit = utils.CALLID_STATUS_PRE_AUDIT
			input.TaskStatusInit = utils.TASK_STATUS_ANNOTATED
		}
	}
	input.NumAudit = numAudit
	totalAnalysts := 0
	userAnalysts := []uint32{}
	numAnal := make(map[uint32]int, 0)
	for _, m := range input.Analysts {
		if m.Num == 0 {
			continue
		}
		totalAnalysts += m.Num
		if _, nok := numAnal[m.UserID]; nok {
			numAnal[m.UserID] += m.Num
			continue
		}
		numAnal[m.UserID] = m.Num
		rel := service.TaskUserRelationItem{
			TID:      taskinfo.Id,
			TaskID:   taskinfo.Taskid,
			UserID:   m.UserID,
			TaskRole: utils.TASK_PAGE_TAB_ANALYS,
			TotalNum: m.Num,
			Status:   utils.TASK_STATUS_ALLOCATED,
		}
		userAnalysts = append(userAnalysts, m.UserID)
		input.Relation = append(input.Relation, rel)
		if input.StatusInit == "" {
			input.StatusInit = utils.TASK_STATUS_PRE_ANALYST
			input.TaskStatusInit = utils.TASK_STATUS_AUDITED
		}
	}
	input.NumAnal = numAnal
	if input.StatusInit == "" {
		seelog.Errorf("DoAllocatTask by taskid, 标注人，审核人，结果分析人至少选一个")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "Select at least one of the annotator, the reviewer, and the result analyst", nil))
		return
	}
	if (totalAnalysts > 0 && totalAnalysts != taskinfo.CallNum) || (totalAuditor > 0 && totalAuditor != taskinfo.CallNum) || (totalAnnot > 0 && totalAnnot != taskinfo.CallNum) {
		seelog.Errorf("DoAllocatTask by taskid, CallNum is no match")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "CallNum is no match", nil))
		return
	}
	//分配任务
	input.Id = taskinfo.Id
	input.TaskAnnotator = utils.UInt32SliceToString(userAnnot)
	input.TaskAuditor = utils.UInt32SliceToString(userAuditor)
	input.TaskAnalysts = utils.UInt32SliceToString(userAnalysts)
	err3 := service.DoAllocatTask(input, "do")
	if err3 != nil {
		seelog.Errorf("DoAllocatTask %s failed :%s", input.TaskId, err3.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "DoAllocatTask task failed", nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 任务管理-预览-审核数量-总数、完成数
func GetTaskAuditStatByTaskId(ctx *gin.Context) {
	seelog.Infof("update one call labelwork detail cmd")
	input := service.GetTaskStatDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("update one call labelwork detail request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)

	if len(input.TaskId) == 0 || userID == 0 {
		seelog.Errorf("update one call labelwork detail request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("DoAllocatTask by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "taskid is error", nil))
		return
	}

	//检查该任务 审核 完成数据
	total, num, _ := service.GetTaskUserDoneNumByTab(input.TaskId, utils.TASK_PAGE_TAB_AUDIT)
	resp := service.GetTaskStatRetDTO{
		Total:    total,
		Finished: num,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// RingType 任务创建
func RingTypeTask(ctx *gin.Context) {
	input := service.RingTypeTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("RingTypeTask request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "userlogin request marsh error", nil))
		return
	}
	if input.Key == "" || input.TaskName == "" || len(input.CallIDs) == 0 || len(input.LableNames) == 0 || input.FinishTime == "" {
		seelog.Errorf("RingTypeTask param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	if len(input.CallIDs) > 500 {
		seelog.Errorf("RingTypeTask callid too more")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid too more, max 500", nil))
		return
	}
	taskKey, ok := utils.RingTypeTaskKeys[input.Key]
	if !ok {
		seelog.Errorf("RingTypeTask key error %s", input.Key)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "key is error", nil))
		return
	}

	task, err := service.CheckTaskName(input.TaskName)
	if err != nil || task != "" {
		seelog.Errorf("taskname existed %s", input.TaskName)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_TASK_NAME_EXIST, "", nil))
		return
	}
	var (
		callids        []string
		dataErr        bool
		labelIDs       []string
		systemLabelID  string
		countryLabelID string
	)
	for _, call := range input.CallIDs {
		if call.CallID == "" || call.URL == "" {
			dataErr = true
			continue
		}
		callids = append(callids, call.CallID)
	}
	for _, labName := range input.LableNames {
		if lab, ok := service.LabelNameInfoCache[labName]; ok {
			if lab.LabelName == "System RingType" {
				systemLabelID = lab.LabelId
			}
			if lab.LabelName == "Country" {
				countryLabelID = lab.LabelId
			}
			labelIDs = append(labelIDs, lab.LabelId)
		} else {
			seelog.Errorf("labName not exist %s", labName)
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "labName not exist", nil))
			return
		}
	}
	if dataErr {
		seelog.Errorf("callids info error or empty")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callids info error or empty", nil))
		return
	}

	creator, err := service.GetUserInfoByLoginName(taskKey.Creator)
	if err != nil || creator.LoginName == "" {
		seelog.Errorf("creator userinfo  error or empty")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "creator userinfo  error or empty", nil))
		return
	}
	taskOperators := []uint32{}
	for _, allocat := range taskKey.Allocator {
		allocator, err := service.GetUserInfoByLoginName(allocat)
		if err != nil || allocator.LoginName == "" {
			continue
		}
		taskOperators = append(taskOperators, allocator.UserID)
	}
	if len(taskOperators) == 0 {
		seelog.Errorf("allocator userinfo  error or empty")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "allocator userinfo  error or empty", nil))
		return
	}
	country := service.GetCountryInfoByUID(creator.UserID)
	finishTime := utils.TimeZone2UTCTime(input.FinishTime, country.ZoneId)
	if finishTime.Before(utils.NowUTC()) {
		seelog.Errorf("create task finishTime less than current time")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "finishTime less than current time", nil))
		return
	}
	req := service.CreateTaskDTO{
		TaskName:      input.TaskName,
		CreatorId:     creator.UserID,
		TaskCreator:   creator.UserName,
		TaskOperators: taskOperators,
		Callid:        callids,
		LabelId:       labelIDs,
	}
	req.StartTimer = utils.NowUTC()
	req.FinishTimer = finishTime

	err, taskid := service.CreateNewTask(req)
	if err != nil {
		seelog.Errorf("create task failed:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task failed", nil))
		return
	}

	err = service.CreateTaskCalls(taskid, req.Callid)
	if err != nil {
		seelog.Errorf("create task calls failed:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task calls failed", nil))
		return
	}
	err = service.CreateRingTypeTaskLabel(input.CallIDs, labelIDs, taskid, systemLabelID, countryLabelID)
	go service.PrepareRingTypeTaskRecordingFiles(taskid, input.CallIDs)

	var ret service.CreateTaskRetDTO
	ret.Taskid = taskid

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

// 任务管理-预览-全部数据
func GetTaskPreviewLabelworkList(ctx *gin.Context) {
	seelog.Infof("get subtask labelwork list cmd")
	input := service.TaskPreviewLabelworkListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get subtask labelwork list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get label list request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)
	if len(input.TaskID) == 0 || userID == 0 {
		seelog.Error("get subtask labelwork list subtask id is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "get label list request marsh error", nil))
		return
	}

	if input.PageNum == 0 {
		seelog.Warnf("get subtask labelwork list request page num is 0")
		input.PageNum = 1
	}
	if input.PageSize == 0 {
		seelog.Warnf("get subtask labelwork list request page size is 0")
		input.PageSize = 100
	}

	seelog.Infof("get subtask labelwork list cmd page num:%d, page size:%d", input.PageNum, input.PageSize)
	taskLabelworkList, total, err := service.GetTaskPreviewCallsList(input.TaskID, "", userID, nil, input.PageNum, input.PageSize)
	if err != nil {
		seelog.Errorf("get subtask labelwork list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get subtask labelwork list failed", nil))
		return
	}
	progress, err := service.GetTaskProgressDetails(ctx, input.TaskID)
	if err != nil {
		seelog.Errorf("get subtask labelwork list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "GetTaskProgressDetails failed", nil))
		return
	}

	taskinfo, err := service.GetTaskInfoByID(input.TaskID)
	if err != nil || taskinfo == nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid is error", nil))
		return
	}

	resp := service.TaskPreviewLabelworkListRetDTO{
		TaskTotal: taskinfo.CallNum,
		Total:     int(total),
		List:      taskLabelworkList,
		Progress:  progress,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}
