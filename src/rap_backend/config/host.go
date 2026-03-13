package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cihub/seelog"
)

var (
	HttpAddr string
	HttpPort int

	gHostCfgPath string

	EsDomain     string
	EsIndex      string
	EsDoc        string
	GlobalConfig HostNodeConfig
)

type HostNodeConfig struct {
	Env            string `json:"env"`
	ServerName     string `json:"servername"`
	HttpServerIp   string `json:"httpserverip"`
	HttpServerPort int    `json:"httpserverport"`
	Desc           string `json:"desc"`
	Enable         bool   `json:"enable"`
	EsDomain       string `json:"es_domain"`
	EsIndex        string `json:"es_index"`
	EsDoc          string `json:"es_doc"`
}

func LoadHostConfigData(data string) int {
	configData := HostNodeConfig{}
	if err := json.Unmarshal([]byte(data), &configData); err != nil {
		seelog.Error("parse node cfg error ", err)
		return -1
	}
	GlobalConfig = configData
	HttpAddr = configData.HttpServerIp
	HttpPort = configData.HttpServerPort
	EsDomain = configData.EsDomain
	EsIndex = configData.EsIndex
	EsDoc = configData.EsDoc

	return 0
}

func LoadHostConfig(cfg string) int {
	gHostCfgPath = cfg

	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		seelog.Error(err)
		return -1
	}

	return LoadHostConfigData(string(file))
}
