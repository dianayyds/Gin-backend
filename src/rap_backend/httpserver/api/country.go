package httpserver

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	"rap_backend/dao"
	"rap_backend/fileprocess"
	"rap_backend/internal/context"
	"rap_backend/service"
	"strings"
)

func CountryList(ctx *gin.Context) {
	input := service.CountryListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("CountryList request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "CountryList request marsh error", nil))
		return
	}

	resp := service.CountryListRetDTO{}

	info, ret := service.CountryList()
	if ret != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "country error", nil))
		return
	}
	resp.List = info
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

// 前期临时入口，方便执行sql
func AdminExec(ctx *gin.Context) {
	input := service.DownloadTaskDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("CountryList request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "CountryList request marsh error", nil))
		return
	}
	userName := context.GetLoginName(ctx)
	if !(userName == "junliang.chen@airudder.com" || userName == "lulu.cao@airudder.com") || input.SubTaskId != "3fa0941b70h84e96b96488ab479b5fe6" || input.TaskId == "" {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "", nil))
		return
	}
	if strings.HasPrefix(input.TaskId, "tsk_") {
		fileprocess.PrepareTaskRecordingFilesLost(input.TaskId)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, input.TaskId, nil))
		return
	}
	resp := service.ResponseRetDTO{}

	err := service.ExecSQL(input.TaskId)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

type AdminFixDataDTO struct {
	TaskId  string   `json:"task_id" form:"task_id"`
	CallIds []string `json:"call_ids" form:"call_ids"`
}

func AdminFixData(ctx *gin.Context) {
	input := AdminFixDataDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("CountryList request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "CountryList request marsh error", nil))
		return
	}

	task, err := dao.GetTaskInfoByID(input.TaskId)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	if len(input.CallIds) == 0 {
		taskCalls, err := dao.GetCallIDsByTaskID(input.TaskId)
		if err != nil {
			ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
			return
		}

		input.CallIds = taskCalls
	}

	err = service.CreateTaskLabel2(service.CreateTaskDTO{
		Callid:  input.CallIds,
		LabelId: strings.Split(task.Labels, ","),
	}, input.TaskId)
	if err != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(nil))
}
