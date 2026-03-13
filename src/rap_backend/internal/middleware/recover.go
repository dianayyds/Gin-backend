package middleware

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				seelog.Errorf("panic err: %+v", err)
				seelog.Errorf(string(debug.Stack()))
				seelog.Error("panic error"+string(debug.Stack()), nil)
			}
		}()
		c.Next()
	}
}
