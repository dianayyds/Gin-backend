package httpserver

import (
	"net/http"
	"rap_backend/internal/context"
	"rap_backend/service"
	"rap_backend/utils"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func SSOLogin(ctx *gin.Context) {
	input := service.SSOLoginDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("ssologin request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "ssologin request marsh error", nil))
		return
	}
	if input.Code == "" {
		seelog.Errorf("code is null")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "code is null", nil))
		return
	}
	referer := ctx.GetHeader("Origin")
	token, err := service.SSOUserLogin(referer, input.Code)
	if err != nil || token == "" {
		seelog.Errorf("SSOUserLogin err:%v", err)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "login failure", nil))
		return
	}
	ret := service.SSOLoginRetDTO{
		Token: token,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

func UserLogin(ctx *gin.Context) {
	input := service.UserLoginDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserLogin request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "userlogin request marsh error", nil))
		return
	}
	if input.LoginName == "" || input.Password == "" {
		seelog.Errorf("UserLogin param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	token, err := service.UserLogin(input.LoginName, input.Password)
	if err != nil || token == "" {
		seelog.Errorf("UserLogin err:%s, name:%s", err, input.LoginName)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "login failure", nil))
		return
	}
	ret := service.SSOLoginRetDTO{
		Token: token,
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}

func UserInfo(ctx *gin.Context) {
	input := service.UserInfoDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserInfoDTO request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "UserInfoDTO request marsh error", nil))
		return
	}
	userID := context.GetUID(ctx)
	if userID == 0 {
		seelog.Errorf("need login")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_USER_NEED_LOGIN, "need login", nil))
		return
	}
	if input.UserID > 0 {
		userID = input.UserID
	}
	resp := service.UserInfoRetDTO{}

	info, ret := service.UserInfo(userID)
	if ret != nil || info == nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "user error", nil))
		return
	}
	resp.Info = *info
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func UserList(ctx *gin.Context) {
	input := service.UserListDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserInfoDTO request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "UserInfoDTO request marsh error", nil))
		return
	}
	ulid := time.Now().UnixNano()
	seelog.Infof("rapUserList_start %d, t:%d", ulid, time.Now().Unix())
	if input.PageNum < 1 {
		input.PageNum = 1
	}
	if input.PageSize < 1 || input.PageSize > 200 {
		input.PageSize = 20
	}
	resp := service.UserListRetDTO{}
	offset := (input.PageNum - 1) * input.PageSize
	info, total, ret := service.UserList(input.LoginName, input.UserType, input.Status, offset, input.PageSize)
	if ret != nil || info == nil {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "user error", nil))
		return
	}
	resp.List = info
	resp.Total = total
	seelog.Infof("rapUserList_end %d, t:%d", ulid, time.Now().Unix())

	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func UserCreate(ctx *gin.Context) {
	input := service.UserCreateDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserCreateDTO request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "UserCreateDTO request marsh error", nil))
		return
	}

	if input.UserName == "" || input.LoginName == "" || input.Password == "" || input.CountryId == 0 || len(input.Roles) == 0 {
		seelog.Errorf("UserCreate param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	if len(input.LoginName) > 30 || len(input.UserName) > 30 {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "name is too long", nil))
		return
	}
	if _, ok := service.CountryCache[uint32(input.CountryId)]; !ok {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "country info error", nil))
		return
	}
	roleInfos, err := service.RoleInfoByIDs(input.Roles)
	if err != nil || len(roleInfos) != len(input.Roles) {
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "role info error", nil))
		return
	}

	//检测密码安全程度
	if utils.CheckPwd(8, 16, input.Password) != nil {
		seelog.Errorf("UserCreate The password must be 8-16 digits and include numbers, uppercase and lowercase letters (%s)", input.Password)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "Password error", nil))
		return
	}
	//时区
	resp := service.ResponseRetDTO{}

	uid, ret := service.CreateUser(input)
	if ret != nil || uid == 0 {
		seelog.Errorf("UserCreate create error:%v", ret)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "user existed", nil))
		return
	}
	resp.Ret = 1
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func UserEdit(ctx *gin.Context) {
	input := service.UserEditDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserEdit request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "UserEdit request marsh error", nil))
		return
	}
	if input.UserID == 0 || input.UserName == "" || input.CountryId == 0 || len(input.Roles) == 0 {
		seelog.Errorf("UserEdit param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	//检测密码安全程度
	if input.Password != "" && utils.CheckPwd(8, 16, input.Password) != nil {
		seelog.Errorf("UserCreate The password must be 8-16 digits and include numbers, uppercase and lowercase letters (%s)", input.Password)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "Password error", nil))
		return
	}
	//时区
	resp := service.ResponseRetDTO{}

	ret := service.EditUser(input)
	if ret != nil {
		seelog.Errorf("UserEdit edit error:%v", ret)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "UserEdit error", nil))
		return
	}
	resp.Ret = 1
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func UserOnOff(ctx *gin.Context) {
	input := service.UserStatusDTO{}
	if err := ctx.Bind(&input); err != nil {
		seelog.Errorf("UserOnOff request marsh error :%s", err.Error())
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_REQUEST, "UserOnOff request marsh error", nil))
		return
	}
	if input.UserID == 0 || input.Status == 0 {
		seelog.Errorf("UserOnOff param error")
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_PARAM, "param error", nil))
		return
	}
	//时区
	resp := service.ResponseRetDTO{}

	ret := service.EditUserStatus(input)
	if ret != nil {
		seelog.Errorf("UserOnOff edit error:%v", ret)
		ctx.JSON(http.StatusOK, NewCommonResp(ERR_SYSTEM_ERROR, "UserOnOff error", nil))
		return
	}
	resp.Ret = 1
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(resp))
}

func CallingIC(ctx *gin.Context) {
	// input := map[string]interface{}{}
	// if err := ctx.Bind(&input); err != nil {
	// 	seelog.Errorf("UserLogin request marsh error :%s", err.Error())
	// 	ctx.JSON(http.StatusOK, NewCommonResp(ERR_COMMON_SINGLE, "userlogin request marsh error", nil))
	// 	return
	// }
	a, e := ctx.GetRawData()
	seelog.Infof("callingic: data: %s, %v", string(a), e)

	// seelog.Infof("callingic: data: %v", input)
	ret := service.SSOLoginRetDTO{
		Token: "token",
	}
	ctx.JSON(http.StatusOK, NewCommonSuccessResp(ret))
}
