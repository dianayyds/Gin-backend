package httpserver

import (
	"net/http"
	"rap_backend/service"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func RoleList(ctx *gin.Context) {
	input := service.RoleListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("RoleList request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "RoleList request marsh error", nil))
		return
	}
	if input.PageNum < 1 {
		input.PageNum = 1
	}
	if input.PageSize < 1 || input.PageSize > 100 {
		input.PageSize = 20
	}

	resp := service.RoleListRetDTO{}

	info, total, ret := service.RoleList((input.PageNum-1)*input.PageSize, input.PageSize)
	if ret != nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "Role error", nil))
		return
	}
	resp.List = info
	resp.Total = total
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}
