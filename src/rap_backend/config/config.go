package config

import (
	"fmt"
	"os"
	"path"
)

var CfgRoot string
var LogFile string
var LogLevel string
var Tostdflag bool

var NodeName string
var DbConfig string
var GaussDbConfig string
var OSSConfig string
var CommonConfig string

var NodeConfig string
var HostConfig string

var BasicFile string
var ServiceRoot string

var gConfig *Config = nil

var RedChache RedisCache

func LoadConfig(configFile string) int {
	if gConfig == nil {
		gConfig = &Config{}
	}

	if gConfig.LoadCofig(configFile) < 0 {
		fmt.Fprintln(os.Stderr, "Read ", configFile, " failed.")
		return -1
	}

	Tostdflag = false
	NodeName = GetNodeName()
	LogLevel = gConfig.GetValue(CFG_LOG_LEVEL)
	if LogLevel == "" {
		LogLevel = "Info"
	}

	if gConfig.GetValue(CFG_LOG_TOSTD) == "true" {
		Tostdflag = true
	}

	CfgRoot = gConfig.GetValue(CFG_CFG_ROOT)
	if CfgRoot == "" {
		fmt.Fprintln(os.Stderr, "No CfgRoot path!")
		return -1
	}

	LogFile = gConfig.GetValue(CFG_LOG_FILE)
	if LogFile == "" {
		fmt.Fprintln(os.Stderr, "No LogFile path!")
		return -1
	}

	DbConfig = "db.config"
	GaussDbConfig = "gaussdb.config"
	HostConfig = NodeName + ".config"
	OSSConfig = "oss.config"
	CommonConfig = "common" + ".config"

	DbConfig = path.Join(CfgRoot, path.Base(DbConfig))
	GaussDbConfig = path.Join(CfgRoot, path.Base(GaussDbConfig))
	HostConfig = path.Join(CfgRoot, path.Base(HostConfig))
	OSSConfig = path.Join(CfgRoot, path.Base(OSSConfig))
	CommonConfig = path.Join(CfgRoot, path.Base(CommonConfig))

	if NodeName == "" {
		fmt.Fprintln(os.Stderr, "No server node name!")
		return -1
	}

	return 0
}

func GetNodeName() string {
	// in config
	nodeName := gConfig.GetValue(CFG_NODENAME)
	if nodeName != "" {
		return nodeName
	}

	return CONTAINER_HEADER_NAME
}
