package context

import (
	"rap_backend/internal/jwtauth"
	"rap_backend/service"

	"github.com/gin-gonic/gin"
)

func GetUID(c *gin.Context) uint32 {
	var uID uint32
	claims, exist := c.Get("claims")
	if !exist {
		return uID
	}
	cla, ok := claims.(*jwtauth.CustomClaims)
	if !ok {
		return uID
	}
	uID = uint32(cla.UserId)
	return uID
}

func GetLoginName(c *gin.Context) string {
	var loginName string
	claims, exist := c.Get("claims")
	if !exist {
		return loginName
	}
	cla, ok := claims.(*jwtauth.CustomClaims)
	if !ok {
		return loginName
	}
	loginName = cla.Username
	return loginName
}

func GetUserInfo(c *gin.Context) *service.UserInfoDetail {
	var uID uint32
	claims, exist := c.Get("claims")
	if !exist {
		return nil
	}
	cla, ok := claims.(*jwtauth.CustomClaims)

	if !ok {
		return nil
	}
	uID = uint32(cla.UserId)
	info, _ := service.UserInfo(uID)
	return info
}
