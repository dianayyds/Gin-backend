package httpserver

import (
	"net/http"
	"rap_backend/config"
	httpserver "rap_backend/httpserver/api"
	"rap_backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
}

func InitRoute() *gin.Engine {
	r := gin.Default()
	rest := Router{}
	//gin.DefaultWriter = seelog.Out
	r.Use(rest.corsMiddleware)
	r.NoRoute(rest.notFound)
	r.NoMethod(rest.noMethod)
	r.Use(middleware.RecoveryMiddleware())
	// r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)))

	apiGroup := r.Group("/rap")
	{
		//登录
		apiGroup.POST("/sso/login", httpserver.SSOLogin)
		apiGroup.POST("/user/login", httpserver.UserLogin)
		apiGroup.POST("/calling/ic", httpserver.CallingIC)
		apiGroup.POST("/task/ringtype", httpserver.RingTypeTask)

	}
	apiGroup.Use(middleware.JWTAuth())
	{

		apiGroup.POST("/task/detail", httpserver.TaskDetail)

		apiGroup.POST("/task/submit/statistics", httpserver.SubmitStatistics)

		//任务创建
		apiGroup.POST("/task/checkcallid", httpserver.CheckUploadCallId)
		apiGroup.POST("/task/label/default", httpserver.LabelDefault)
		apiGroup.POST("/task/create", httpserver.CreateTaskV2)
		//任务管理
		apiGroup.POST("/task/getlist", httpserver.GetTaskList)
		apiGroup.POST("/task/download", httpserver.DownloadTask)
		apiGroup.POST("/task/delete", httpserver.DeleteTask)

		apiGroup.POST("/task/statics", httpserver.GetTaskAuditStatByTaskId)

		//预览-全部数据
		apiGroup.POST("/task/preview", httpserver.GetTaskPreviewLabelworkList)
		// apiGroup.POST("/labelwork/getsubtaskdetaillist", httpserver.GetSubTaskLabelworkList)

		//任务分配
		apiGroup.POST("/task/allocat/list", httpserver.GetAllocatTaskList)
		apiGroup.POST("/task/allocat/do", httpserver.DoAllocatTask)
		apiGroup.POST("/task/allocat/redo", httpserver.ReDoAllocatTask)
		apiGroup.POST("/task/allocat/redo-single", httpserver.RedoAllocatTaskSingle)

		//通话标注
		apiGroup.POST("/task/annotat/list", httpserver.GetAnnotatTaskList)
		//标注 callids
		apiGroup.POST("/task/annotat/calls", httpserver.GetAnnotatTaskCallsByTaskId)
		// apiGroup.POST("/task/getdetailbytaskid", httpserver.GetTaskDetailByTaskId)
		//完成数量
		apiGroup.POST("/task/annotat/statics", httpserver.GetAnnotatTaskStatByTaskId)
		// apiGroup.POST("/task/getsubtaskstatics", httpserver.GetSubTaskStatics)
		//标注callid 详情
		apiGroup.POST("/task/annotat/callid/detail", httpserver.GetOneCallLabelWorkDetail)
		// apiGroup.POST("/labelwork/getonedetail", httpserver.GetOneCallLabelWorkDetail)
		//标注callid 内容标注
		apiGroup.POST("/task/annotat/label/update", httpserver.UpdateOneCallLabelWorkDetail)
		// apiGroup.POST("/labelwork/updateonedetail", httpserver.UpdateOneCallLabelWorkDetail)
		//完成任务标注-提交
		apiGroup.POST("/task/annotat/label/done", httpserver.TaskLabelWorkDone)
		//任务标注-预览
		apiGroup.POST("/task/annotat/preview", httpserver.GetAnnotatTaskPreviewLabelworkList)

		//通话审核 任务列表
		apiGroup.POST("/task/audit/list", httpserver.GetAuditTaskList)
		//审核 callids
		apiGroup.POST("/task/audit/calls", httpserver.GetAuditTaskCallsByTaskId)
		//完成数量
		apiGroup.POST("/task/audit/statics", httpserver.GetAuditTaskStatByTaskId)
		//审核callid 详情
		apiGroup.POST("/task/audit/callid/detail", httpserver.GetAuditOneCallLabelWorkDetail)
		//审核callid 内容审核
		apiGroup.POST("/task/audit/label/update", httpserver.UpdateAuditOneCallLabelWorkDetail)
		//审核callid - 驳回
		apiGroup.POST("/task/audit/label/reject", httpserver.AuditReject)
		//完成任务审核-提交
		apiGroup.POST("/task/audit/label/done", httpserver.TaskAuditLabelWorkDone)
		//通话审核-预览
		apiGroup.POST("/task/audit/preview", httpserver.GetAuditTaskPreviewLabelworkList)

		//结果分析 任务列表
		apiGroup.POST("/task/analyst/list", httpserver.GetAnalystTaskList)
		//分析 callids
		apiGroup.POST("/task/analyst/calls", httpserver.GetAnalystTaskCallsByTaskId)
		//完成数量
		apiGroup.POST("/task/analyst/statics", httpserver.GetAnalystTaskStatByTaskId)
		//结果分析callid 详情
		apiGroup.POST("/task/analyst/callid/detail", httpserver.GetAnalystOneCallLabelWorkDetail)
		//分析callid 内容分析
		apiGroup.POST("/task/analyst/label/update", httpserver.UpdateAnalystOneCallLabelWorkDetail)
		//完成任务分析-提交
		apiGroup.POST("/task/analyst/label/done", httpserver.TaskAnalystLabelWorkDone)
		//结果分析-预览
		apiGroup.POST("/task/analyst/preview", httpserver.GetAnalystTaskPreviewLabelworkList)

		//字段管理
		apiGroup.POST("/label/getlist", httpserver.GetLabelList)
		apiGroup.POST("/label/create", httpserver.CreateLabel)
		apiGroup.POST("/label/edit", httpserver.EditLabel)
		apiGroup.POST("/label/del", httpserver.DelLabel)
		//字段模版管理
		apiGroup.GET("/label/template/getList", httpserver.GetLabelTemplateList)
		apiGroup.POST("/label/template/create", httpserver.CreateLabelTemplate)
		apiGroup.POST("/label/template/update", httpserver.UpdateLabelTemplateById)
		apiGroup.POST("/label/template/delete", httpserver.DeleteLabelTemplate)

		//用户管理
		apiGroup.POST("/user/info", httpserver.UserInfo)
		apiGroup.POST("/user/list", httpserver.UserList)
		apiGroup.POST("/user/create", httpserver.UserCreate)
		apiGroup.POST("/user/edit", httpserver.UserEdit)
		apiGroup.POST("/user/onoff", httpserver.UserOnOff)
		apiGroup.POST("/country/list", httpserver.CountryList)
		apiGroup.POST("/role/list", httpserver.RoleList)

		//执行sql
		apiGroup.POST("/admin/exec", httpserver.AdminExec)

		apiGroup.POST("/admin/fixdata", httpserver.AdminFixData)

	}
	return r
}

func (r Router) notFound(q *gin.Context) {
	//q.JSON(http.StatusNotFound, tools.Common{Code: -1, Message: "failed", Body: "request target not found"})
	return
}

func (r Router) noMethod(q *gin.Context) {
	//q.JSON(http.StatusMethodNotAllowed, tools.Common{Code: -1, Message: "failed", Body: "request method doesn't support"})
	return
}

func (r Router) corsMiddleware(q *gin.Context) {
	q.Writer.Header().Set("Access-Control-Allow-Origin", q.GetHeader("Origin"))
	q.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	q.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Token")
	q.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if q.Request.Method == "OPTIONS" {
		q.AbortWithStatus(200)
		return
	}
	q.Next()
}

func (r Router) tokenVerify(q *gin.Context) {

	if q.Request.Method == "OPTIONS" {
		q.AbortWithStatus(200)
		return
	}
	token := q.GetHeader("Authorization")
	if token != "REPLACE_WITH_INTERNAL_AUTH_TOKEN" {
		q.AbortWithStatus(http.StatusForbidden)
		return
	}
	q.Next()
}

func (r Router) purviewVerify(q *gin.Context) {

	if q.Request.Method == "OPTIONS" {
		q.AbortWithStatus(200)
		return
	}
	uri := q.Request.RequestURI
	pre, ok := config.PURVIEWAPI[uri]
	if !ok {
		pre = ""
		// q.AbortWithStatus(http.StatusUnauthorized)
		// return
	}
	if pre != "" {
	}

	q.Next()
}
