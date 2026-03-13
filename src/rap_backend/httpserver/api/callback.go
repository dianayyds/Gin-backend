package httpserver

//import (
//	"github.com/gin-gonic/gin"
//	"gitlab.airudder.com/airudder-production/sms_gateway/internal"
//	"gitlab.airudder.com/airudder-production/sms_gateway/internal/log"
//	"gitlab.airudder.com/airudder-production/sms_gateway/services/sms"
//	"net/http"
//	"strings"
//)
//
//func CallbackHandlerPost(ctx *gin.Context) {
//	ispsKey := strings.TrimSpace(ctx.Param("ispskey"))
//	log.Infof(ctx, "接收到来自%s的请求，requestUrl:%s", ispsKey, ctx.Request.RequestURI)
//	err := sms.GetManager().HanderCallback(ispsKey, ctx)
//	if err != nil {
//		log.Warnf(ctx, "回执处理失败")
//		ctx.JSON(http.StatusInternalServerError,
//			internal.NewCommonResp(internal.ERR_COMMON_SINGLE,
//				internal.AccessAlertMessage(internal.ERR_COMMON_SINGLE),
//				nil))
//		return
//	}
//	ctx.JSON(http.StatusOK,
//		internal.NewCommonResp(internal.NORMAL_SINGLE,
//			internal.AccessAlertMessage(internal.NORMAL_SINGLE),
//			nil))
//	return
//}
