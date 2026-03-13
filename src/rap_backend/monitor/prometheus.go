package monitor

import (
	"os"
	"rap_backend/config"
	"rap_backend/version"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	SERVICE = version.SERVICE
	CBGROUP = os.Getenv("CB_GROUP")
	Buckets = func() []float64 {
		return []float64{0.001, 0.002, 0.005, 0.01,
			0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.8, 1, 2, 5, 10}
	}

	pktloss90Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "pktloss90",
			Help:      "90 line for pktloss.",
		},

		[]string{"node", "linename"},
	)

	pktloss95Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "pktloss95",
			Help:      "90 line for pktloss.",
		},

		[]string{"node", "linename"},
	)

	audioFileErrorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "audioFileError",
			Help:      "number of when audio file is error.",
		},

		[]string{"node", "audioFileName"},
	)

	icallStreamErrorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "icallStreamErr",
			Help:      "Create icall stream err.",
		},

		[]string{"node", "icallAddress"},
	)

	sendIcallStreamErrorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "sendIcallStreamErr",
			Help:      "Send Icall stream err.",
		},

		[]string{"node", "icallAddress"},
	)

	fcsStreamErrorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "fcsStreamErr",
			Help:      "Create fcs stream err.",
		},

		[]string{"node", "fcsAddress"},
	)

	portPoolSizeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "portPoolSize",
			Help:      "port pool free size",
		},

		[]string{"node"},
	)

	lossrateLessGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "lossrateLessGauge",
			Help:      "number of loss 3 > rate > 0 now.",
		},

		[]string{"node", "mediaip"},
	)

	lossrate3Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "lossrate3Gauge",
			Help:      "number of loss rate > 3 now.",
		},

		[]string{"node", "mediaip"},
	)

	lossrate5Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "lossrate5Gauge",
			Help:      "number of loss 10 < rate < 20 now.",
		},

		[]string{"node", "mediaip"},
	)

	lossrate10Gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "lossrate10Gauge",
			Help:      "number of loss 3 < rate < 10 now.",
		},

		[]string{"node", "mediaip"},
	)

	lossrateBigGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: CBGROUP,
			Subsystem: SERVICE,
			Name:      "lossrateBigGauge",
			Help:      "number of loss rate > 20 now.",
		},

		[]string{"node", "mediaip"},
	)
)

func PrometheusPktLoss90(linename string, v float32) {
	node := config.GetNodeName()
	pktloss90Gauge.With(
		prometheus.Labels{
			"node":     node,
			"linename": linename,
		},
	).Set(float64(v))
}

func PrometheusPktLoss95(linename string, v float32) {
	node := config.GetNodeName()
	pktloss95Gauge.With(
		prometheus.Labels{
			"node":     node,
			"linename": linename,
		},
	).Set(float64(v))
}

func PrometheusPortPoolSize(v int64) {
	node := config.GetNodeName()
	portPoolSizeGauge.With(
		prometheus.Labels{
			"node": node,
		},
	).Set(float64(v))
}

func PrometheusAudioFileError(errAudioFileName string, v int64) {
	node := config.GetNodeName()
	audioFileErrorGauge.With(
		prometheus.Labels{
			"node":          node,
			"audioFileName": errAudioFileName,
		},
	).Add(float64(v))
}

func PrometheusIcallStreamError(errIcallAddress string, v int64) {
	node := config.GetNodeName()
	icallStreamErrorGauge.With(
		prometheus.Labels{
			"node":         node,
			"icallAddress": errIcallAddress,
		},
	).Add(float64(v))
}

func PrometheusFCSStreamError(errFCSAddress string, v int64) {
	node := config.GetNodeName()
	fcsStreamErrorGauge.With(
		prometheus.Labels{
			"node":       node,
			"fcsAddress": errFCSAddress,
		},
	).Add(float64(v))
}

func PrometheusSendIcallStreamError(errIcallAddress string, v int64) {
	node := config.GetNodeName()
	sendIcallStreamErrorGauge.With(
		prometheus.Labels{
			"node":         node,
			"icallAddress": errIcallAddress,
		},
	).Add(float64(v))
}

func PrometheusLossRateLess(mediaIp string, v int64) {
	node := config.GetNodeName()
	lossrateLessGauge.With(
		prometheus.Labels{
			"node":    node,
			"mediaip": mediaIp,
		},
	).Add(float64(v))
}

func PrometheusLossRate3(mediaIp string, v int64) {
	node := config.GetNodeName()
	lossrate3Gauge.With(
		prometheus.Labels{
			"node":    node,
			"mediaip": mediaIp,
		},
	).Add(float64(v))
}

func PrometheusLossRate5(mediaIp string, v int64) {
	node := config.GetNodeName()
	lossrate5Gauge.With(
		prometheus.Labels{
			"node":    node,
			"mediaip": mediaIp,
		},
	).Add(float64(v))
}

func PrometheusLossRate10(mediaIp string, v int64) {
	node := config.GetNodeName()
	lossrate10Gauge.With(
		prometheus.Labels{
			"node":    node,
			"mediaip": mediaIp,
		},
	).Add(float64(v))
}

func PrometheusLossRateBig(mediaIp string, v int64) {
	node := config.GetNodeName()
	lossrateBigGauge.With(
		prometheus.Labels{
			"node":    node,
			"mediaip": mediaIp,
		},
	).Add(float64(v))
}

func init() {
	prometheus.MustRegister(portPoolSizeGauge)
	prometheus.MustRegister(lossrate3Gauge)
	prometheus.MustRegister(lossrate5Gauge)
	prometheus.MustRegister(lossrate10Gauge)
	prometheus.MustRegister(lossrateBigGauge)
	prometheus.MustRegister(lossrateLessGauge)
	prometheus.MustRegister(audioFileErrorGauge)
	prometheus.MustRegister(icallStreamErrorGauge)
	prometheus.MustRegister(fcsStreamErrorGauge)
	prometheus.MustRegister(pktloss90Gauge)
	prometheus.MustRegister(pktloss95Gauge)
	prometheus.MustRegister(sendIcallStreamErrorGauge)
}
