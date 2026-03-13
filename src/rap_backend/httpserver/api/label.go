package httpserver

import (
	"net/http"
	"rap_backend/config"
	"rap_backend/internal/context"
	"rap_backend/service"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func GetLabelList(ctx *gin.Context) {
	seelog.Infof("get label list cmd")
	input := service.GetLabelDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get label list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get label list request marsh error", nil))
		return
	}
	if input.PageNum <= 0 {
		seelog.Warnf("get label list request page num is 0")
		input.PageNum = 1
	}
	if input.PageSize <= 0 {
		seelog.Warnf("get label list request page size is 0")
		input.PageSize = 100
	}
	orderBy := "is_editable,label_name"
	if input.PageTab == "label_page" {
		orderBy = "is_editable desc,update_time desc"
	}
	seelog.Infof("get label list cmd page num:%d, page size:%d", input.PageNum, input.PageSize)
	labelInfoList, cnt, err := service.GetLabelList(input.LabelName, input.PageNum, input.PageSize, config.USER_STATUS_NORMAL, orderBy)
	if err != nil {
		seelog.Errorf("get label list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get label list error failed", nil))
		return
	}
	var ret = service.GetLabelRetDTO{}
	ret.LabelCnt = cnt
	ret.LabelInfoList = labelInfoList
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

func CreateLabel(ctx *gin.Context) {
	seelog.Infof("receive create label cmd")
	input := service.CreateLabelDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("create label request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "check callid request marsh error", nil))
		return
	}
	if input.LabelName == "" {
		seelog.Errorf("CreateLabel param error:%v", input)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	input.LabelCreator = context.GetLoginName(ctx)

	_, err := service.CreateNewLabel(input)
	if err != nil {
		seelog.Errorf("create label failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	go service.RefreshLabelInfosLocalCache()
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func EditLabel(ctx *gin.Context) {
	seelog.Infof("receive create label cmd")
	input := service.EditLabelDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("create label request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "check callid request marsh error", nil))
		return
	}
	if input.LabelName == "" || input.LabelId == "" {
		seelog.Errorf("CreateLabel param error:%v", input)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}

	err := service.EditLabel(input)
	if err != nil {
		seelog.Errorf("edit label failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func DelLabel(ctx *gin.Context) {
	seelog.Infof("receive create label cmd")
	input := service.DelLabelDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("del label request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "check callid request marsh error", nil))
		return
	}
	if input.LabelId == "" {
		seelog.Errorf("del Label param error:%v", input)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}

	err := service.DelLabel(input)
	if err != nil {
		seelog.Errorf("edit label failed :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	resp := service.ResponseRetDTO{
		Ret: 1,
	}

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}
