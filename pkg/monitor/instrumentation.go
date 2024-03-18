package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	EventCount       prometheus.Gauge
	Info             *prometheus.GaugeVec
	PromHttpRespTime *prometheus.HistogramVec
}

func NewMetrics(promCamRegistry prometheus.Registerer) *Metrics {
	m := &Metrics{
		EventCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices.",
		}),
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "info",
			Help:      "Information about the Event Stats environment.",
		},
			[]string{"version"}),
		PromHttpRespTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "vsapi_http_response_time",
			Help:      "Duration of HTTP requests.",
			// Buckets:   prometheus.LinearBuckets(0.6, 3, 15),
			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		},
			[]string{"path", "status"}),
	}
	promCamRegistry.MustRegister(
		m.EventCount,
		m.Info,
		m.PromHttpRespTime,
	)
	return m
}
