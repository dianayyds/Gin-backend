package httpserver

import (
	"fmt"
	"net/http"
	"rap_backend/internal/context"
	"rap_backend/service"
	"rap_backend/utils"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

// 结果分析页面-任务了列表.
func GetAnalystTaskList(ctx *gin.Context) {
	seelog.Infof("get GetAnalystTaskList cmd")
	input := service.GetTaskListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get GetAnalystTaskList request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get GetAuditTaskList request marsh error", nil))
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
		PageTab: utils.TASK_PAGE_TAB_ANALYS,
		StatusList: []string{
			utils.TASK_STATUS_ALLOCATED,
			utils.TASK_STATUS_ANNOTATING,
			utils.TASK_STATUS_ANNOTATED,
			utils.TASK_STATUS_AUDITING,
			utils.TASK_STATUS_AUDITED,
			utils.TASK_STATUS_ANALYSTING,
			utils.TASK_STATUS_PRE_ANALYST,
			utils.TASK_STATUS_ANNOTATED,
			utils.TASK_STATUS_COMPLETED,
		},
		UserIDs:     []uint32{userinfo.UserID},
		SearchType:  input.SearchType,
		SearchValue: input.SearchValue,
		PageNum:     input.PageNum,
		PageSize:    input.PageSize,
	}
	if userinfo.DataAll { //全局数据，查看所有
		req.UserIDs = []uint32{}
	}
	taskInfoList, cnt, err2 := service.GetTaskListByFilter(req)
	if err2 != nil {
		seelog.Errorf("get task list error failed")
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

// 结果页面-获取任务下callid列表.
func GetAnalystTaskCallsByTaskId(ctx *gin.Context) {
	seelog.Infof("receive get task  detail by taskid cmd")
	input := service.GetTaskCallIdListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get task detail by taskid request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get task detail by taskid request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)
	var lastID int32
	//var status string
	if userID == 0 {
		seelog.Errorf("get task detail by taskid, need login")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "need login", nil))
		return
	}
	if input.TaskId == "" {
		seelog.Errorf("get task detail by taskid, taskid is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task id is null", nil))
		return
	}

	if input.PageSize == 0 {
		seelog.Warnf("get task detail request page size is 0")
		input.PageSize = 10
	}
	//获取未标注数据
	status := []string{
		utils.TASK_STATUS_PRE_ANALYST, //待结果分析
		utils.TASK_STATUS_ANALYSTED,   //结果分析完成
	}
	if input.NoMark == 1 {
		input.PageNum = 1
		status = []string{
			utils.TASK_STATUS_PRE_ANALYST, //待结果分析
		}
	}
	if input.PageNum == 0 {
		seelog.Warnf("get task detail request page num is 0")
		input.PageNum = 1
		//第一次进入标注页面，先定位上次退出是哪个位置，哪个页面
		lastID = service.GetTaskUserRelationLastID(input.TaskId, userID, utils.TASK_PAGE_TAB_ANALYS)
		if lastID > 0 {
			lastPage := service.GetTaskUserLastPage(input.TaskId, utils.TASK_PAGE_TAB_ANALYS, userID, lastID, input.PageSize)
			if lastPage > 0 {
				input.PageNum = lastPage
			}
		}
	}
	fmt.Println(lastID)

	calls, loc, total, err := service.GetTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_ANALYS,
		userID,
		status,
		input.PageNum,
		input.PageSize,
		lastID,
	)
	if err != nil {
		seelog.Errorf("get task detail by taskid failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "create task failed", nil))
		return
	}
	resp := service.GetTaskCallIdListRetDTO{
		TaskId:       input.TaskId,
		CallIdList:   calls,
		Total:        total,
		LastLocation: loc,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话审核页面-数量统计-总数、完成数
func GetAnalystTaskStatByTaskId(ctx *gin.Context) {
	seelog.Infof("get GetAnalystTaskStatByTaskId cmd")
	input := service.GetTaskStatDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get GetAnalystTaskStatByTaskId request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)

	if len(input.TaskId) == 0 || userID == 0 {
		seelog.Errorf("get GetAnalystTaskStatByTaskId request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("GetAnalystTaskStatByTaskId by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "taskid is error", nil))
		return
	}

	//检查该用户负责的callid数量，及完成数据
	total, num := service.GetTaskStatusNumList(
		input.TaskId,
		[]string{utils.TASK_STATUS_ANALYSTED},
		utils.TASK_PAGE_TAB_ANALYS,
		userID,
	)
	resp := service.GetTaskStatRetDTO{
		Total:    total,
		Finished: num,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 结果分析-获取某个callid下 label info.
func GetAnalystOneCallLabelWorkDetail(ctx *gin.Context) {
	seelog.Infof("get one call labelwork detail cmd")
	input := service.GetOneCallLabelWorkDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get one call labelwork detail request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}

	if len(input.CallId) == 0 {
		seelog.Errorf("get one call labelwork detail request callid is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid is null", nil))
		return
	}

	if len(input.TaskId) == 0 {
		seelog.Errorf("get one call labelwork detail request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}

	LabelWorkDetail := service.GetOneCallLabelWork(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_ANALYS)
	recordURL := service.GetRecordURLByTaskIDCallID(input.TaskId, input.CallId)
	creator, annor, auditor := service.GetOneCallLabelWorkUsers(input.TaskId, input.CallId)
	resp := service.GetOneCallLabelWorkRetDTO{
		CallId:           input.CallId,
		CallRecordURL:    recordURL,
		LabelWorkInfoLst: LabelWorkDetail,
		Creator:          creator,
		Annotator:        annor,
		Auditor:          auditor,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 结果分析-审核某个callid-保存内容.
func UpdateAnalystOneCallLabelWorkDetail(ctx *gin.Context) {
	seelog.Infof("update one call labelwork detail cmd")
	input := service.UpdateOneCallLabelWorkDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("update one call labelwork detail request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)

	if len(input.CallId) == 0 {
		seelog.Errorf("update one call labelwork detail request callid is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid is null", nil))
		return
	}

	if len(input.TaskId) == 0 {
		seelog.Errorf("update one call labelwork detail request SubTaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "SubTaskId is null", nil))
		return
	}
	taskCallId, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || taskCallId.CallId == "" {
		seelog.Errorf("GetTaskCallInfoByTaskIDCallID empty or err:%v, taskID:%s, callid:%s", err, input.TaskId, input.CallId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "task call error", nil))
		return
	}
	if taskCallId.MultiAnalysts != userID {
		seelog.Errorf("GetTaskCallInfoByTaskIDCallID not match :taskID:%v, taskID:%s, callid:%s", err, input.TaskId, input.CallId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot analyze the results of this task", nil))
		return
	}

	if taskCallId.Status != utils.TASK_STATUS_ANALYSTED &&
		taskCallId.Status != utils.TASK_STATUS_PRE_ANALYST {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid status error", nil))
		return
	}

	//保持数据分析内容
	err = service.UpdateAnalystOneCallLabelWork(input.LabelWorkInfoLst)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	upd := map[string]interface{}{
		"actual_analysts": userID,
		"status":          utils.TASK_STATUS_ANALYSTED,
	}
	err = service.UpdateTaskCallByWhere(input.TaskId, input.CallId, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByWhere err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//更新任务表状态
	service.UpdateTaskStatusByTaskStatus(input.TaskId, utils.TASK_STATUS_ANALYSTING, utils.TASK_STATUS_AUDITED)
	//taskuserrelation记录最后一次操作的id,及状态
	err = service.UpdateTaskUserRelationForLabel(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_ANALYS, userID, taskCallId.SerialNumber)
	if err != nil {
		seelog.Errorf("UpdateTaskUserRelationForLabel err:%s", err.Error())
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 结果分析-任务完成-全部提交
func TaskAnalystLabelWorkDone(ctx *gin.Context) {
	seelog.Infof("update one call labelwork detail cmd")
	input := service.TaskLabelWorkDoneDTO{}
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
	if taskinfo.TaskStatus != utils.TASK_STATUS_ANALYSTING {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid status error", nil))
		return
	}

	taskUserRel, err := service.GetTaskUserRelationInfo(input.TaskId, utils.TASK_PAGE_TAB_ANALYS, userID)
	if err != nil || taskUserRel.TID == 0 {
		seelog.Errorf("GetTaskUserRelationInfo empty or err:%v, taskID:%s", err, input.TaskId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot analyze the results of this task", nil))
		return
	}
	//检查该用户负责的callid数量，及完成数据
	total, num := service.GetTaskStatusNumList(
		input.TaskId,
		[]string{utils.TASK_STATUS_ANALYSTED},
		utils.TASK_PAGE_TAB_ANALYS,
		userID,
	)
	if total != num {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "还有callid未审核", nil))
		return
	}

	analystedCallIds, _, err := service.GetOriginTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_ANALYS,
		userID,
		[]string{utils.TASK_STATUS_ANALYSTED},
		1,
		1000000,
	)
	if len(analystedCallIds) == 0 || err != nil {
		seelog.Errorf("DoAllocatTask by taskid, taskid status error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "review at least one", nil))
		return
	}

	err = service.UpdateCallIdsStatus(analystedCallIds, utils.TASK_STATUS_COMPLETED)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	upd := map[string]interface{}{
		"status":   utils.TASK_STATUS_ANALYSTED,
		"done_num": num,
	}
	err = service.UpdateTaskUserRelationByWhere(input.TaskId, userID, utils.TASK_PAGE_TAB_ANALYS, upd)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//检查整个任务的状态，判断任务是否进入下一个状态
	isDone, err := service.GetTaskCurrentStatusIsDone(
		input.TaskId,
		utils.TASK_PAGE_TAB_ANALYS,
		utils.TASK_STATUS_ANALYSTED,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//任务当前阶段全部完成 ，任务下一个状态是什么
	if isDone {
		nextStatus, err := service.GetTaskNextStatus(input.TaskId)
		if err != nil {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
			return
		}
		if nextStatus != "" {
			err = service.UpdateTaskStatusByTaskID(input.TaskId, nextStatus)
			if err != nil {
				ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
				return
			}
		}
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 结果分析-预览-个人部分
func GetAnalystTaskPreviewLabelworkList(ctx *gin.Context) {
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
	taskLabelworkList, total, err := service.GetTaskPreviewCallsList(input.TaskID, utils.TASK_PAGE_TAB_ANALYS, userID, nil, input.PageNum, input.PageSize)
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
