package httpserver

import (
	"errors"
)

const (
	NORMAL_SINGLE      = 0
	ERR_COMMON_REQUEST = -1 //Receive request data error
	ERR_COMMON_PARAM   = -2 //parameter error, please check
	ERR_SYSTEM_ERROR   = -3 //System error, please try again
	ERR_COMMON_TIMEOUT = -4 //Request time out
	ERR_UNKNOWN_ERROR  = -5 //unknown error

	ERR_USER_NEED_LOGIN    = -101 //need login
	ERR_USER_TOKEN_INVALID = -102 //invalid token

	ERR_TASK_NAME_EXIST = -201 //任务名已存在

	// Body Central error definition area, 20000-29999
	NsqBindConsumerError = "20004"
)

var ErrAlerts = map[int]string{
	NORMAL_SINGLE:          "success",
	ERR_COMMON_REQUEST:     "Receive request data error",
	ERR_COMMON_PARAM:       "Parameter error, please check",
	ERR_SYSTEM_ERROR:       "System error, please try again",
	ERR_COMMON_TIMEOUT:     "Request time out",
	ERR_UNKNOWN_ERROR:      "Unknown error",
	ERR_USER_NEED_LOGIN:    "Need login",
	ERR_USER_TOKEN_INVALID: "Invalid token",
	ERR_TASK_NAME_EXIST:    "Task name already exists",
}

type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

type CommonListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Body    struct {
		Total int64       `json:"total"`
		List  interface{} `json:"list"`
	} `json:"body"`
}

func NewCommonResp(code int, msg string, params interface{}) CommonResp {
	if msg == "" {
		msg = AccessAlertMessage(code)
	}
	return CommonResp{
		code,
		msg,
		params,
	}
}

func NewCommonSuccessResp(params interface{}) CommonResp {
	return CommonResp{
		0,
		"success",
		params,
	}
}
func NewCommonListResp(code int, msg string, totalCount int64, listObjects interface{}) CommonListResp {
	clp := CommonListResp{
		Code:    code,
		Message: msg,
	}
	clp.Body.Total = totalCount
	clp.Body.List = listObjects
	return clp
}

func AccessAlertMessage(code int) string {
	if message, ok := ErrAlerts[code]; ok {
		return message
	}
	return ""
}

func NewError(errno int) error {
	return errors.New(ErrAlerts[errno])
}
