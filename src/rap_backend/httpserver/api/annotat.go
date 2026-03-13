package httpserver

import (
	"net/http"
	"rap_backend/internal/context"
	"rap_backend/service"
	"rap_backend/utils"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

// 通话标注页面-任务列表.
func GetAnnotatTaskList(ctx *gin.Context) {
	seelog.Infof("get GetAnnotatTaskList cmd")
	input := service.GetTaskListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get GetAnnotatTaskList request marsh error :%s", err.Error())
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
		PageTab:     utils.TASK_PAGE_TAB_ANNOTAT,
		StatusList:  []string{utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATING},
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

// 通话标注页面-获取任务下callid列表.
func GetAnnotatTaskCallsByTaskId(ctx *gin.Context) {
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
	status := []string{utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATED}
	if input.NoMark == 1 {
		input.PageNum = 1
		status = []string{utils.TASK_STATUS_ALLOCATED}
	}
	if input.PageNum == 0 {
		seelog.Warnf("get task detail request page num is 0")
		input.PageNum = 1
		//第一次进入标注页面，先定位上次退出是哪个位置，哪个页面
		lastID = service.GetTaskUserRelationLastID(input.TaskId, userID, utils.TASK_PAGE_TAB_ANNOTAT)
		if lastID > 0 {
			lastPage := service.GetTaskUserLastPage(input.TaskId, utils.TASK_PAGE_TAB_ANNOTAT, userID, lastID, input.PageSize)
			if lastPage > 0 {
				input.PageNum = lastPage
			}
		}
	}

	calls, loc, total, err := service.GetTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_ANNOTAT,
		userID,
		status,
		input.PageNum,
		input.PageSize,
		lastID,
	)
	if err != nil {
		seelog.Errorf("get task detail by taskid failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "GetTaskCallsList error", nil))
		return
	}
	resp := service.GetTaskCallIdListRetDTO{
		TaskId:       input.TaskId,
		CallIdList:   calls,
		Total:        total,
		LastLocation: loc,
	}

	if input.NoMark == 1 {
		_, num := service.GetTaskStatusNumList(
			input.TaskId,
			[]string{utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATED},
			utils.TASK_PAGE_TAB_ANNOTAT,
			userID,
		)
		resp.Statistics.Total = int64(num)
		resp.Statistics.Done = int64(num) - total
	} else {
		_, num := service.GetTaskStatusNumList(
			input.TaskId,
			[]string{utils.TASK_STATUS_ANNOTATED},
			utils.TASK_PAGE_TAB_ANNOTAT,
			userID,
		)
		resp.Statistics.Total = total
		resp.Statistics.Done = int64(num)
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话标注任务-数量统计-总数、完成数
func GetAnnotatTaskStatByTaskId(ctx *gin.Context) {
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

	//检查该用户负责的callid数量，及完成数据
	total, num := service.GetTaskStatusNumList(
		input.TaskId,
		utils.GetCallIdRemainStatus(utils.TASK_STATUS_ANNOTATED),
		utils.TASK_PAGE_TAB_ANNOTAT,
		userID,
	)
	resp := service.GetTaskStatRetDTO{
		Total:    total,
		Finished: num,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话标注-预览-个人部分
func GetAnnotatTaskPreviewLabelworkList(ctx *gin.Context) {
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
	taskLabelworkList, total, err := service.GetTaskPreviewCallsList(input.TaskID, utils.TASK_PAGE_TAB_ANNOTAT, userID, nil, input.PageNum, input.PageSize)
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
