package db

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/cihub/seelog"
	"gorm.io/gorm"
)

var CFG_DB string
var CFG_USR string
var CFG_PWD string
var CFG_HOST string
var CFG_PROTO string
var CFG_PORT string
var CFG_MAXIDLE int
var CFG_MAXCONN int
var CFG_IDLETIMEOUT int
var CFG_ENABLE bool

var GMysalDB *gorm.DB

type DataBaseCfg struct {
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

func GetTxDB(ctx context.Context, txDB *gorm.DB) *gorm.DB {
	if txDB != nil {
		return txDB
	}

	return GMysalDB
}

func LoadConfigData(data string) int {
	configData := DataBaseCfg{}
	if err := json.Unmarshal([]byte(data), &configData); err != nil {
		seelog.Error("parse db cfg error ", err)
		return -1
	}

	CFG_DB = configData.Db
	if CFG_DB == "" {
		seelog.Error("Get db error")
		return -1
	}

	CFG_USR = configData.Usr
	if CFG_USR == "" {
		seelog.Error("Get usr error")
		return -1
	}

	CFG_PWD = configData.Pwd
	//if CFG_PWD == "" {
	//	seelog.Error("Get pwd error")
	//	return -1
	//}

	CFG_HOST = configData.Host
	if CFG_HOST == "" {
		seelog.Error("Get host error")
		return -1
	}

	CFG_PORT = configData.Port
	if CFG_PORT == "" {
		seelog.Error("Get port error")
		return -1
	}

	CFG_MAXIDLE = configData.MaxIdle
	if CFG_MAXIDLE == 0 {
		seelog.Error("Get maxidle error")
		return -1
	}

	CFG_MAXCONN = configData.MaxConn
	if CFG_MAXCONN == 0 {
		seelog.Error("Get maxconn error")
		return -1
	}

	CFG_IDLETIMEOUT = configData.IdleTimeout
	if CFG_IDLETIMEOUT == 0 {
		seelog.Error("Get idle timeout error")
		return -1
	}

	CFG_PROTO = configData.Proto
	if CFG_PROTO != "tcp" && CFG_PROTO != "udp" {
		seelog.Error("Get proto error")
		return -1
	}

	CFG_ENABLE = configData.Enable
	return 0
}

func LoadConfig(cfg string) int {
	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		seelog.Error(err)
		return -1
	}

	return LoadConfigData(string(file))
}

func StopMysqlServer() {
	if GMysalDB != nil {

	}
}

func StartMysqlServer() int {
	var err error
	GMysalDB, err = InitMysql(CFG_USR, CFG_PWD, CFG_HOST, CFG_PORT, CFG_DB, CFG_MAXIDLE, CFG_MAXCONN, CFG_IDLETIMEOUT)
	if err != nil || GMysalDB.Error != nil {
		seelog.Errorf("start mysql server error:%s", err.Error())
		return -1
	}
	seelog.Infof("start mysql server successfully:%s", CFG_HOST)
	return 0
}
