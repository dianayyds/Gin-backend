package httpserver

import (
	"gorm.io/gorm"
	"net/http"
	"rap_backend/internal/context"
	"rap_backend/service"
	"rap_backend/utils"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

// 通话审核页面-任务了列表.
func GetAuditTaskList(ctx *gin.Context) {
	input := service.GetTaskListDTO{}
	if err := ctx.Bind(&input); err != nil {
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
		PageTab:     utils.TASK_PAGE_TAB_AUDIT,
		StatusList:  []string{utils.TASK_STATUS_ALLOCATED, utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_ANNOTATING, utils.TASK_STATUS_AUDITING},
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

// 通话审核页面-获取任务下callid列表.
func GetAuditTaskCallsByTaskId(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task id is null", nil))
		return
	}

	if input.PageSize == 0 {
		seelog.Warnf("get task detail request page size is 0")
		input.PageSize = 10
	}
	//获取未标注数据
	status := []string{
		utils.CALLID_STATUS_PRE_AUDIT, //待审核
		utils.TASK_STATUS_AUDITED,     //已审核
	}
	if input.NoMark == 1 {
		input.PageNum = 1
		status = []string{
			utils.CALLID_STATUS_PRE_AUDIT, //待审核
		}
	}
	if input.PageNum == 0 {
		seelog.Warnf("get task detail request page num is 0")
		input.PageNum = 1
		//第一次进入标注页面，先定位上次退出是哪个位置，哪个页面
		lastID = service.GetTaskUserRelationLastID(input.TaskId, userID, utils.TASK_PAGE_TAB_AUDIT)
		if lastID > 0 {
			lastPage := service.GetTaskUserLastPage(input.TaskId, utils.TASK_PAGE_TAB_AUDIT, userID, lastID, input.PageSize)
			if lastPage > 0 {
				input.PageNum = lastPage
			}
		}
	}

	calls, loc, total, err := service.GetTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_AUDIT,
		userID,
		status,
		input.PageNum,
		input.PageSize,
		lastID,
	)
	if err != nil {
		seelog.Errorf("get task detail by taskid failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "create task failed", nil))
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
			[]string{
				utils.CALLID_STATUS_PRE_AUDIT, //待审核
				utils.TASK_STATUS_AUDITED,     //已审核
			},
			utils.TASK_PAGE_TAB_AUDIT,
			userID,
		)
		resp.Statistics.Total = int64(num)
		resp.Statistics.Done = int64(num) - total
	} else {
		_, num := service.GetTaskStatusNumList(
			input.TaskId,
			[]string{utils.TASK_STATUS_AUDITED},
			utils.TASK_PAGE_TAB_AUDIT,
			userID,
		)
		resp.Statistics.Total = total
		resp.Statistics.Done = int64(num)
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话审核页面-数量统计-总数、完成数
func GetAuditTaskStatByTaskId(ctx *gin.Context) {
	seelog.Infof("get GetAuditTaskStatByTaskId cmd")
	input := service.GetTaskStatDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get GetAuditTaskStatByTaskId request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)

	if len(input.TaskId) == 0 || userID == 0 {
		seelog.Errorf("get GetAuditTaskStatByTaskId request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("GetAuditTaskStatByTaskId by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid is error", nil))
		return
	}

	//检查该用户负责的callid数量，及完成数据
	total, num := service.GetTaskStatusNumList(
		input.TaskId,
		utils.GetCallIdRemainStatus(utils.TASK_STATUS_AUDITED),
		utils.TASK_PAGE_TAB_AUDIT,
		userID,
	)
	resp := service.GetTaskStatRetDTO{
		Total:    total,
		Finished: num,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话审核-获取某个callid下 label info.
func GetAuditOneCallLabelWorkDetail(ctx *gin.Context) {
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

	LabelWorkDetail := service.GetOneCallLabelWork(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_AUDIT)
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

// 通话审核-审核某个callid-保存内容.
func UpdateAuditOneCallLabelWorkDetail(ctx *gin.Context) {
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
	taskCall, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || taskCall.CallId == "" {
		seelog.Errorf("GetTaskCallInfoByTaskIDCallID empty or err:%v, taskID:%s, callid:%s", err, input.TaskId, input.CallId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task call error", nil))
		return
	}
	if taskCall.MultiAuditor != userID {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot review this task", nil))
		return
	}

	//保持了下审核信息
	err = service.UpdateAuditOneCallLabelWork(input.LabelWorkInfoLst)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//status 限制
	if taskCall.Status != utils.CALLID_STATUS_PRE_AUDIT &&
		taskCall.Status != utils.TASK_STATUS_AUDITED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "call status error", nil))
		return
	}

	upd := map[string]interface{}{
		"actual_auditor": userID,
		"status":         utils.TASK_STATUS_AUDITED,
	}

	err = service.UpdateTaskCallByWhere(input.TaskId, input.CallId, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByWhere err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//审核 - 每次都尝试更新任务表状态
	err = service.UpdateTaskStatusByTaskStatus(input.TaskId, utils.TASK_STATUS_AUDITING, utils.TASK_STATUS_ANNOTATED)
	if err != nil {
		seelog.Errorf("UpdateTaskStatusByTaskStatus err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//taskuserrelation记录最后一次操作的id,及状态
	err = service.UpdateTaskUserRelationForLabel(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_AUDIT, userID, taskCall.SerialNumber)
	if err != nil {
		seelog.Errorf("UpdateTaskUserRelationForLabel err:%s", err.Error())
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话审核-任务完成-全部提交
func TaskAuditLabelWorkDone(ctx *gin.Context) {
	seelog.Infof("update one call labelwork detail cmd")
	input := service.TaskLabelWorkDoneDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("update one call labelwork detail request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)

	if len(input.TaskId) == 0 || userID == 0 || input.Total == 0 {
		seelog.Errorf("update one call labelwork detail request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil {
		seelog.Errorf("DoAllocatTask by taskid, taskid is error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid is error", nil))
		return
	}

	//状态限制： 必须先保存一条审核，才能提交
	//查询 至少有一条已审核callid

	hasAuditcallids, _, err := service.GetOriginTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_AUDIT,
		userID,
		[]string{utils.TASK_STATUS_AUDITED},
		1,
		1000000,
	)
	if len(hasAuditcallids) == 0 {
		seelog.Errorf("DoAllocatTask by taskid, taskid status error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "review at least one", nil))
		return
	}
	taskUserRel, err := service.GetTaskUserRelationInfo(input.TaskId, utils.TASK_PAGE_TAB_AUDIT, userID)
	if err != nil || taskUserRel.TID == 0 {
		seelog.Errorf("GetTaskUserRelationInfo empty or err:%v, taskID:%s", err, input.TaskId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "You cannot review this task", nil))
		return
	}

	//查询已审核的callid,更新为待数据分析 ， 未审核部分提交了也更新？
	callids, _, err := service.GetOriginTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_AUDIT,
		userID,
		[]string{utils.TASK_STATUS_AUDITED, utils.CALLID_STATUS_PRE_AUDIT},
		1,
		1000000,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	if len(callids) == 0 {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "不存在的callid", nil))
		return
	}

	// 确认待审核、已审核的num == 前端传的num， 防止审核中，又有callid流转过来
	if len(callids) != input.Total {
		ctx.JSON(http.StatusOK, NewCommonResp(-6, "callid数据不匹配，请刷新页面", nil))
		return
	}

	callIdNextStatus, err := service.GetCallIDNextStatus(taskinfo, utils.TASK_STATUS_AUDITED)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	err = service.UpdateCallIdsStatus(callids, callIdNextStatus)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//计算task callid是否全部审核完成，  ->  计算最新的task status

	//1. 计算user task callid 是否全部审核完成 -> 更新taskuserrelation

	userLeftCallids, err := service.GetUserTaskCallsByStatus(
		input.TaskId,
		userID,
		utils.TASK_STATUS_AUDITED,
		[]string{
			utils.TASK_STATUS_ALLOCATED,   //已分配，待标注
			utils.TASK_STATUS_ANNOTATED,   //已标注
			utils.CALLID_STATUS_PRE_AUDIT, //待审核
			utils.TASK_STATUS_AUDITED,     //已审核

		},
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	upd := map[string]interface{}{
		"done_num": gorm.Expr("done_num + ?", len(callids)),
	}

	if len(userLeftCallids) == 0 {
		upd["status"] = utils.TASK_STATUS_AUDITED
	}

	err = service.UpdateTaskUserRelationByWhere(input.TaskId, userID, utils.TASK_PAGE_TAB_AUDIT, upd)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//2 .计算task callid 是否全部标注完成
	TaskLeftCallids, err := service.GetUserTaskCallsByStatus(
		input.TaskId,
		0,
		"",
		[]string{
			utils.TASK_STATUS_ALLOCATED,   //已分配，待标注
			utils.TASK_STATUS_ANNOTATED,   //已标注
			utils.CALLID_STATUS_PRE_AUDIT, //待审核
			utils.TASK_STATUS_AUDITED,     //已审核
		},
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	if len(TaskLeftCallids) == 0 {
		//最后一次部分标注提交 -> task结束 -> 状态流转
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

// 通话审核-预览-个人部分
func GetAuditTaskPreviewLabelworkList(ctx *gin.Context) {
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
	taskLabelworkList, total, err := service.GetTaskPreviewCallsList(input.TaskID, utils.TASK_PAGE_TAB_AUDIT, userID, nil, input.PageNum, input.PageSize)
	if err != nil {
		seelog.Errorf("get preview labelwork list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get preview labelwork list failed", nil))
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

func auditRejectAll(ctx *gin.Context, input service.UpdateRejectOneCallLabelWorkDTO, userID uint32) {
	//先验证任务信息
	taskinfo, err := service.GetTaskInfoByID(input.TaskId)
	if err != nil || taskinfo == nil || !utils.IsInUint32Slice(taskinfo.Auditor, userID) {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid or userId is error", nil))
		return
	}

	//获取任意一个callid info
	taskCall, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || taskCall.CallId == "" {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task call error", nil))
		return
	}
	if taskCall.MultiAuditor != userID {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot review this task", nil))
		return
	}

	//获取用户负责的所有, 处于可审核、待审核的callid, 都要驳回
	TaskCallids, err := service.GetUserTaskCallsByStatus(
		input.TaskId,
		taskCall.MultiAnnotator, //todo zb这里的实际标记人，不一定是现在的标记人。 （比如zb先标注，然后重新分配给gll）， 会导致全部驳回，只驳回zb的标注
		utils.TASK_PAGE_TAB_ANNOTAT,
		[]string{
			utils.CALLID_STATUS_PRE_AUDIT, //待审核
			utils.TASK_STATUS_AUDITED,     //已审核
		},
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	var TCids []int32
	for _, callid := range TaskCallids {
		TCids = append(TCids, callid.TcId)
	}

	upd := map[string]interface{}{
		"reject_reason":  input.RejectReason,
		"actual_auditor": userID,
		"status":         utils.TASK_STATUS_ALLOCATED,
	}

	count, err := service.UpdateTaskCallByTcIDs(TCids, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByWhere err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//更新任务表状态
	err = service.UpdateTaskStatusByTaskIDOnly(input.TaskId, utils.TASK_STATUS_ANNOTATING)
	if err != nil {
		seelog.Errorf("UpdateTaskStatusByTaskIDOnly err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//taskuserrelation 更新
	relationUpd := map[string]interface{}{
		"status":   utils.TASK_STATUS_ANNOTATING,
		"done_num": gorm.Expr("done_num - ?", len(TaskCallids)),
	}

	err = service.UpdateTaskUserRelationByWhere(
		input.TaskId,
		taskCall.MultiAnnotator,
		utils.TASK_PAGE_TAB_ANNOTAT,
		relationUpd,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	resp := service.ResponseRetDTO{
		Ret:   1,
		Count: count,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))

}

func AuditReject(ctx *gin.Context) {
	input := service.UpdateRejectOneCallLabelWorkDTO{}
	if err := ctx.Bind(&input); err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get one call labelwork detail request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)
	if len(input.CallId) == 0 && !input.IsAll {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid is null", nil))
		return
	}

	if len(input.TaskId) == 0 || input.RejectReason == "" {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "RejectReason is null", nil))
		return
	}

	if input.IsAll {
		auditRejectAll(ctx, input, userID)
		return
	}

	taskCall, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || taskCall.CallId == "" {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "task call error", nil))
		return
	}
	if taskCall.MultiAuditor != userID {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot review this task", nil))
		return
	}

	//status 限制
	if taskCall.Status != utils.CALLID_STATUS_PRE_AUDIT && taskCall.Status != utils.TASK_STATUS_AUDITED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "call status error", nil))
		return
	}

	upd := map[string]interface{}{
		"reject_reason":  input.RejectReason,
		"actual_auditor": userID,
		"status":         utils.TASK_STATUS_ALLOCATED,
	}

	err = service.UpdateTaskCallByWhere(input.TaskId, input.CallId, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByWhere err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//更新任务表状态
	err = service.UpdateTaskStatusByTaskIDOnly(input.TaskId, utils.TASK_STATUS_ANNOTATING)
	if err != nil {
		seelog.Errorf("UpdateTaskStatusByTaskIDOnly err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	//taskuserrelation 更新
	relationUpd := map[string]interface{}{
		"status":   utils.TASK_STATUS_ANNOTATING,
		"done_num": gorm.Expr("done_num - ?", 1),
	}

	err = service.UpdateTaskUserRelationByWhere(
		input.TaskId,
		taskCall.ActualAnnotator,
		utils.TASK_PAGE_TAB_ANNOTAT,
		relationUpd,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	resp := service.ResponseRetDTO{
		Ret:   1,
		Count: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}
