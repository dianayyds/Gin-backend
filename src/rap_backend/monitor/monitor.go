package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gRun bool
var gSeq uint64
var gArg chan string

func StopMonitor() {
	gRun = false
}

func StartMonitor() int {
	go startMonitor()
	return 0
}

func startMonitor() {
	gSeq = 0
	gRun = true
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	//seelog.Info("montior start listen the port:" + config.GNodeConfig.MonitorAddr)
	//seelog.Error(http.ListenAndServe(config.GNodeConfig.MonitorAddr, nil))
}
