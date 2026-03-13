package middleware

import (
	"net/http"
	"rap_backend/config"
	"rap_backend/internal/context"
	"rap_backend/internal/jwtauth"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth ...
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"message": "请求未携带token，无权限访问",
				"body":    nil,
			})
			c.Abort()
			return
		}
		if !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"message": "token无效，无权限访问",
				"body":    nil,
			})
			c.Abort()
			return
		}
		token = strings.Replace(token, "Bearer ", "", -1)
		j := jwtauth.NewJWT()
		claims, err := j.ParserToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -102,
				"message": err.Error(),
				"body":    nil,
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		userInfo := context.GetUserInfo(c)
		if userInfo == nil || userInfo.Status == config.USER_STATUS_FORBID {
			c.JSON(http.StatusOK, gin.H{
				"code":    -103,
				"message": "user status error",
				"body":    nil,
			})
			c.Abort()
			return
		}
	}
}
