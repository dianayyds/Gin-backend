package version

import (
	"os"
	"strings"
)

var (
	VERSION   = "1.0.0"
	SERVICE   = GetenvService()
	USERAGENT = SERVICE + " " + VERSION
)

func getenvService() string {
	s := os.Getenv("SERVICE")
	a := strings.Split(s, ".")
	if len(a) > 0 {
		return a[len(a)-1]
	}

	return os.Getenv("SERVICENAME")
}

func GetenvService() string {
	return strings.Replace(
		getenvService(), "-", "", -1)
}
