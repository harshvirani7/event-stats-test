package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	EventCount               prometheus.Gauge
	Info                     *prometheus.GaugeVec
	DurationCountByEventType *prometheus.HistogramVec
	DurationCountByCameraId  *prometheus.HistogramVec
	SuccessRequest           prometheus.Gauge
	ErrorRequest             prometheus.Gauge
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
		DurationCountByEventType: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "request_duration_seconds",
			Help:      "Duration of the request.",
			// 4 times larger for apdex score
			// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
			// Buckets: prometheus.LinearBuckets(0.1, 5, 5),
			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
		DurationCountByCameraId: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "myapp",
			Name:      "request_duration_seconds_totalEventCountByCameraId",
			Help:      "Duration of the request totalEventCountByCameraId",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"status", "method"}),
		SuccessRequest: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "success_request_count",
			Help:      "Total no of success requests.",
		}),
		ErrorRequest: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "error_request_count",
			Help:      "Total no of error requests.",
		}),
	}
	promCamRegistry.MustRegister(m.EventCount, m.Info, m.DurationCountByEventType, m.DurationCountByCameraId, m.SuccessRequest, m.ErrorRequest)
	return m
}
