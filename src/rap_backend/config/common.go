package config

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"io/ioutil"
)

var (
	BIGDATA_HOST  string
	DMROBOT_HOST  string
	DMROBOT_TOKEN string
)

type CommonNodeConfig struct {
	BigdataHost  string `json:"bigdata_host"`
	DmrobotHost  string `json:"dmrobot_host"`
	DmrobotToken string `json:"dmrobot_token"`
}

func LoadCommonConfigData(data string) int {
	configData := CommonNodeConfig{}
	if err := json.Unmarshal([]byte(data), &configData); err != nil {
		seelog.Error("parse common cfg error ", err)
		return -1
	}
	BIGDATA_HOST = configData.BigdataHost
	DMROBOT_HOST = configData.DmrobotHost
	DMROBOT_TOKEN = configData.DmrobotToken

	return 0
}

func LoadCommonConfig(cfg string) int {
	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		seelog.Error(err)
		return -1
	}

	return LoadCommonConfigData(string(file))
}
