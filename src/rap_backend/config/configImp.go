package config

import (
	"io/ioutil"
	"strings"
)

type Config struct {
	kvMap map[string]string
}

func (this *Config) LoadCofig(configFile string) int {
	file, er := ioutil.ReadFile(configFile)
	if er != nil {
		return -1
	}

	if this.kvMap == nil {
		this.kvMap = make(map[string]string)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		params := strings.Split(line, "=")
		if len(params) != 2 {
			continue
		}

		key := strings.TrimSpace(params[0])
		value := strings.TrimSpace(params[1])

		this.kvMap[key] = value
	}

	return 0
}

func (this *Config) GetValue(key string) string {
	return this.kvMap[key]
}
