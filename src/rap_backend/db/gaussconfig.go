package db

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cihub/seelog"
	"gorm.io/gorm"
)

var CFG_GAUSS_DB string
var CFG_GAUSS_USR string
var CFG_GAUSS_PWD string
var CFG_GAUSS_HOST string
var CFG_GAUSS_PROTO string
var CFG_GAUSS_PORT string
var CFG_GAUSS_MAXIDLE int
var CFG_GAUSS_MAXCONN int
var CFG_GAUSS_IDLETIMEOUT int
var CFG_GAUSS_ENABLE bool

var GGaussMysalDB *gorm.DB

type GaussDataBaseCfg struct {
	Db          string `json:"db"`
	Usr         string `json:"usr"`
	Pwd         string `json:"pwd"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Proto       string `json:"proto"`
	MaxIdle     int    `json:"maxidle"`
	MaxConn     int    `json:"maxconn"`
	IdleTimeout int    `json:"idletimeout"`
	Enable      bool   `json:"enable"`
}

func LoadGaussConfigData(data string) int {
	configData := GaussDataBaseCfg{}
	if err := json.Unmarshal([]byte(data), &configData); err != nil {
		seelog.Error("parse gauss db cfg error ", err)
		return -1
	}

	CFG_GAUSS_DB = configData.Db
	if CFG_GAUSS_DB == "" {
		seelog.Error("Get gauss db error")
		return -1
	}

	CFG_GAUSS_USR = configData.Usr
	if CFG_GAUSS_USR == "" {
		seelog.Error("Get gauss usr error")
		return -1
	}

	CFG_GAUSS_PWD = configData.Pwd
	//if CFG_PWD == "" {
	//	seelog.Error("Get pwd error")
	//	return -1
	//}

	CFG_GAUSS_HOST = configData.Host
	if CFG_GAUSS_HOST == "" {
		seelog.Error("Get gauss host error")
		return -1
	}

	CFG_GAUSS_PORT = configData.Port
	if CFG_GAUSS_PORT == "" {
		seelog.Error("Get gauss port error")
		return -1
	}

	CFG_GAUSS_MAXIDLE = configData.MaxIdle
	if CFG_GAUSS_MAXIDLE == 0 {
		seelog.Error("Get gauss maxidle error")
		return -1
	}

	CFG_GAUSS_MAXCONN = configData.MaxConn
	if CFG_GAUSS_MAXCONN == 0 {
		seelog.Error("Get gauss maxconn error")
		return -1
	}

	CFG_GAUSS_IDLETIMEOUT = configData.IdleTimeout
	if CFG_GAUSS_IDLETIMEOUT == 0 {
		seelog.Error("Get gauss idle timeout error")
		return -1
	}

	CFG_GAUSS_PROTO = configData.Proto
	if CFG_GAUSS_PROTO != "tcp" && CFG_GAUSS_PROTO != "udp" {
		seelog.Error("Get gauss proto error")
		return -1
	}

	CFG_GAUSS_ENABLE = configData.Enable
	return 0
}

func LoadGaussConfig(cfg string) int {
	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		seelog.Error(err)
		return -1
	}

	return LoadGaussConfigData(string(file))
}

func StopGaussMysqlServer() {
	if GGaussMysalDB != nil {
	}
}

func StartGaussMysqlServer() int {
	var err error
	GGaussMysalDB, err = InitMysql(CFG_GAUSS_USR, CFG_GAUSS_PWD, CFG_GAUSS_HOST, CFG_GAUSS_PORT, CFG_GAUSS_DB, CFG_GAUSS_MAXIDLE, CFG_GAUSS_MAXCONN, CFG_GAUSS_IDLETIMEOUT)
	if err != nil || GGaussMysalDB.Error != nil {
		seelog.Errorf("start gauss mysql server error:%s", err.Error())
		return -1
	}
	seelog.Infof("start gauss mysql server successfully:%s", CFG_HOST)
	return 0
}
