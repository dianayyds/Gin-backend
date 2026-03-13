package httpserver

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
	ltservice "rap_backend/service/label_template"
)

func GetLabelTemplateList(ctx *gin.Context) {
	input := ltservice.GetLabelTemplateReq{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("get label template list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "get label template list request marsh error", nil))
		return
	}
	resp, err := ltservice.GetLabelTemplateSvc().GetLabelTemplateList(input)
	if err != nil {
		seelog.Errorf("get label template list error failed, err: %+v", err)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "get label template list error failed"+err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func CreateLabelTemplate(ctx *gin.Context) {
	input := ltservice.CreateLabelTemplateReq{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("create label template list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "create label template list request marsh error", nil))
		return
	}
	err := ltservice.GetLabelTemplateSvc().CreateLabelTemplate(input)
	if err != nil {
		seelog.Errorf("create label template list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(nil))
}

func UpdateLabelTemplateById(ctx *gin.Context) {
	input := ltservice.UpdateLabelTemplateReq{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("update label template list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "update label template list request marsh error", nil))
		return
	}
	err := ltservice.GetLabelTemplateSvc().UpdateLabelTemplateById(input.ID, input)
	if err != nil {
		seelog.Errorf("update label template list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "update label template list error failed"+err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(nil))
}

func DeleteLabelTemplate(ctx *gin.Context) {
	input := ltservice.DeleteLabelTemplateReq{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("delete label template list request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "delete label template list request marsh error", nil))
		return
	}
	err := ltservice.GetLabelTemplateSvc().DeleteLabelTemplate(input)
	if err != nil {
		seelog.Errorf("delete label template list error failed")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "delete label template list error failed"+err.Error(), nil))
		return
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(nil))
}
