package config

const (
	CFG_NODENAME       = "NodeName"
	CFG_CFG_ROOT       = "CfgRoot"
	CFG_LOG_FILE       = "LogFile"
	CFG_LOG_TOSTD      = "LogToSTD"
	CFG_LOG_LEVEL      = "LogLevel"
	MaxMsgQueueSize    = 1024
	TABLE_CDR_PKTLOSS  = "cdr_pktloss"
	TABLE_CDR_ANSWERED = "cdr_postbrother_"

	USER_STATUS_NORMAL = 1
	USER_STATUS_FORBID = 2
)

var PURVIEWAPI = map[string]string{
	// "/rap/user/info": "user_manage_list",
}
