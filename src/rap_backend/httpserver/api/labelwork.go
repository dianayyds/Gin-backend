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

// func GetSubTaskLabelworkList(ctx *gin.Context) {
// 	seelog.Infof("get subtask labelwork list cmd")
// 	input := service.GetSubtaskLabelworkListDTO{}
// 	if err := ctx.Bind(&input); err != nil {
// 		seelog.Errorf("get subtask labelwork list request marsh error :%s", err.Error())
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "get label list request marsh error", nil))
// 		return
// 	}
// 	if len(input.SubtaskIdList) == 0 {
// 		seelog.Error("get subtask labelwork list subtask id is null")
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "get label list request marsh error", nil))
// 		return
// 	}

// 	if input.PageNum == 0 {
// 		seelog.Warnf("get subtask labelwork list request page num is 0")
// 		input.PageNum = 1
// 	}
// 	if input.PageSize == 0 {
// 		seelog.Warnf("get subtask labelwork list request page size is 0")
// 		input.PageSize = 100
// 	}

// 	seelog.Infof("get subtask labelwork list cmd page num:%d, page size:%d", input.PageNum, input.PageSize)
// 	subtaskLabelworkList, err := service.GetSubtaskLabelworkList(input.SubtaskIdList, input.PageNum, input.PageSize)
// 	if err != nil {
// 		seelog.Errorf("get subtask labelwork list error failed")
// 		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "get subtask labelwork list failed", nil))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, NewCommonSuccessResp(subtaskLabelworkList))
// }

// 获取某个callid下 label info.
func GetOneCallLabelWorkDetail(ctx *gin.Context) {
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

	item, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || item.CallId == "" {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "task call info empty", nil))
		return
	}

	LabelWorkDetail := service.GetOneCallLabelWork(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_ANNOTAT)
	recordURL := service.GetRecordURLByTaskIDCallID(input.TaskId, input.CallId)
	creator, annor, auditor := service.GetOneCallLabelWorkUsers(input.TaskId, input.CallId)
	resp := service.GetOneCallLabelWorkRetDTO{
		CallId:           input.CallId,
		CallRecordURL:    recordURL,
		LabelWorkInfoLst: LabelWorkDetail,
		Creator:          creator,
		Annotator:        annor,
		Auditor:          auditor,
		RejectReason:     item.RejectReason,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话标注-标注某个callid-保存内容.
func UpdateOneCallLabelWorkDetail(ctx *gin.Context) {
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
		seelog.Errorf("update one call labelwork detail request TaskId is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "TaskId is null", nil))
		return
	}

	taskCallId, err := service.GetTaskCallInfoByTaskIDCallID(input.TaskId, input.CallId)
	if err != nil || taskCallId.CallId == "" {
		seelog.Errorf("GetTaskCallInfoByTaskIDCallID empty or err:%v, taskID:%s, callid:%s", err, input.TaskId, input.CallId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "task call error", nil))
		return
	}
	if taskCallId.MultiAnnotator != userID {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "You cannot mark this task", nil))
		return
	}

	if taskCallId.Status != utils.TASK_STATUS_ALLOCATED &&
		taskCallId.Status != utils.TASK_STATUS_ANNOTATED {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "callid status error", nil))
		return
	}

	//保存标记内容
	err = service.UpdateOneCallLabelWork(input.LabelWorkInfoLst)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	upd := map[string]interface{}{
		"actual_annotator": userID,
		"status":           utils.TASK_STATUS_ANNOTATED,
		"reject_reason":    "",
	}

	//更新callid 状态为已标注
	err = service.UpdateTaskCallByWhere(input.TaskId, input.CallId, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByWhere err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
	}
	//更新任务表状态 task: 已分配 -> 标注中
	err = service.UpdateTaskStatusByTaskStatus(input.TaskId, utils.TASK_STATUS_ANNOTATING, utils.TASK_STATUS_ALLOCATED)
	if err != nil {
		seelog.Errorf("UpdateTaskStatusByTaskStatus err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
	}

	//taskuserrelation记录最后一次操作的id,及状态
	err = service.UpdateTaskUserRelationForLabel(input.TaskId, input.CallId, utils.TASK_PAGE_TAB_ANNOTAT, userID, taskCallId.SerialNumber)
	if err != nil {
		seelog.Errorf("UpdateTaskUserRelationForLabel err:%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 通话标注-任务完成-全部提交
func TaskLabelWorkDone(ctx *gin.Context) {
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

	//task 状态验证
	//if taskinfo.TaskStatus != utils.TASK_STATUS_ANNOTATING {
	//	seelog.Errorf("DoAllocatTask by taskid, taskid status error")
	//	ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "taskid status error", nil))
	//	return
	//}
	//taskid - userid - role  => 权限控制
	taskUserRel, err := service.GetTaskUserRelationInfo(input.TaskId, utils.TASK_PAGE_TAB_ANNOTAT, userID)
	if err != nil || taskUserRel.TID == 0 {
		seelog.Errorf("GetTaskUserRelationInfo empty or err:%v, taskID:%s", err, input.TaskId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "You cannot mark this task", nil))
		return
	}

	//查询已标注的callid,更新为待审核
	callids, _, err := service.GetOriginTaskCallsList(
		input.TaskId,
		utils.TASK_PAGE_TAB_ANNOTAT,
		userID,
		[]string{utils.TASK_STATUS_ANNOTATED},
		1,
		1000000,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	if len(callids) == 0 {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "没有需要提交的callid", nil))
		return
	}

	callIdNextStatus, err := service.GetCallIDNextStatus(taskinfo, utils.TASK_STATUS_ANNOTATED)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	err = service.UpdateCallIdsStatus(
		callids,
		callIdNextStatus,
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//1. 计算user task callid 是否全部标注完成 -> 更新taskuserrelation
	userLeftAnnotatedCallids, err := service.GetUserTaskCallsByStatus(
		input.TaskId,
		userID,
		utils.TASK_PAGE_TAB_ANNOTAT,
		[]string{utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_ALLOCATED},
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	upd := map[string]interface{}{
		"done_num": gorm.Expr("done_num + ?", len(callids)),
	}

	if len(userLeftAnnotatedCallids) == 0 {
		upd["status"] = utils.TASK_STATUS_ANNOTATED
	}

	err = service.UpdateTaskUserRelationByWhere(input.TaskId, userID, utils.TASK_PAGE_TAB_ANNOTAT, upd)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	//2 .计算task callid 是否全部标注完成 ->  计算最新的task status

	TaskLeftAnnotatedCallids, err := service.GetUserTaskCallsByStatus(
		input.TaskId,
		0,
		"",
		[]string{utils.TASK_STATUS_ANNOTATED, utils.TASK_STATUS_ALLOCATED},
	)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	if len(TaskLeftAnnotatedCallids) == 0 {
		//最后一次部分标注提交 -> task结束
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
