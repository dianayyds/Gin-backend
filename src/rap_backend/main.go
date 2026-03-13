package main

import (
	"fmt"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"rap_backend/config"
	"rap_backend/crontab"
	"rap_backend/db"
	"rap_backend/fileprocess"
	"rap_backend/httpserver"
	"rap_backend/log"
	"rap_backend/utils"

	"github.com/cihub/seelog"

	"syscall"
)

func init() {
	seelog.RegisterCustomFormatter("RapLogFormat", log.CreateRapLogFormatter)
}

func args() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "No argrument")
		return -1
	}

	argv := os.Args[1]
	if argv == "debug" {
		//debug.SetDebug(true)
		return 0

	}

	config.BasicFile = argv
	if config.BasicFile == "" {
		fmt.Fprintln(os.Stderr, "No basicconfig file")
		return -1
	}

	return 0
}

func main() {
	if args() < 0 {
		return
	}

	if config.LoadConfig(config.BasicFile) < 0 {
		fmt.Fprintln(os.Stderr, "Load configure file error ", config.BasicFile)
		return
	}

	// init log
	if log.InitLog(config.LogLevel,
		config.LogFile, config.Tostdflag) < 0 {
		fmt.Fprintln(os.Stderr, "Init log file error ", config.LogFile)
		return
	}

	defer seelog.Flush()

	seelog.Infof("Welcome to %s !", config.GetNodeName())
	if db.LoadConfig(config.DbConfig) < 0 {
		seelog.Error("Load db config error ", config.DbConfig)
		return
	}

	if db.LoadGaussConfig(config.GaussDbConfig) < 0 {
		seelog.Error("Load db config error ", config.DbConfig)
		return
	}

	if config.LoadHostConfig(config.HostConfig) < 0 {
		seelog.Error("Load Host config error ", config.HostConfig)
		return
	}

	if fileprocess.LoadConfig(config.OSSConfig) < 0 {
		seelog.Error("Load Oss config error ", config.OSSConfig)
		return
	}

	if config.LoadCommonConfig(config.CommonConfig) < 0 {
		seelog.Error("load Common config error ", config.CommonConfig)
		return
	}

	err := fileprocess.InitOSS()
	if err != nil {
		seelog.Error("init OSS error ", config.OSSConfig)
		return
	}
	utils.Env = config.GlobalConfig.Env
	if db.StartMysqlServer() < 0 {
		seelog.Error("start mysql server error")
		return
	}

	defer db.StopMysqlServer()

	if db.StartGaussMysqlServer() < 0 {
		seelog.Error("start mysql gauss server error")
		return
	}

	defer db.StopGaussMysqlServer()

	// monitor.StartMonitor()
	// defer monitor.StopMonitor()

	crontab.CronInit()
	httpserver.NewHttpServer(config.HttpPort).Start()

	seelog.Info("rap_backend start successfully")

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	// soft kill
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	signal.Stop(signalChan)
	seelog.Infof("See you next time at %s !", config.GetNodeName())
}
