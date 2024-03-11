package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	EventCount                          prometheus.Gauge
	Info                                *prometheus.GaugeVec
	DurationCountByEventType            *prometheus.HistogramVec
	DurationCountByCameraId             *prometheus.HistogramVec
	SuccessRequest                      prometheus.Gauge
	ErrorRequest                        prometheus.Gauge
	StoreEventDataSuccess               prometheus.Gauge
	TotalEventCountByTypeSuccess        prometheus.Gauge
	TotalEventCountByCameraIdSuccess    prometheus.Gauge
	EventCountSummaryByCameraIdSuccess  prometheus.Gauge
	EventCountSummaryByEventTypeSuccess prometheus.Gauge
	SummaryByCameraIdSuccess            prometheus.Gauge
	SummaryByEventTypeSuccess           prometheus.Gauge
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
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3},
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
		StoreEventDataSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "store_event_data_success_count",
			Help:      "Total no of StoreEventData success requests.",
		}),
		TotalEventCountByTypeSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "total_event_count_by_type_success_count",
			Help:      "Total no of TotalEventCountByType success requests.",
		}),
		TotalEventCountByCameraIdSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "total_event_count_by_camera_id_success_count",
			Help:      "Total no of TotalEventCountByCameraId success requests.",
		}),
		EventCountSummaryByCameraIdSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "event_count_summary_by_camera_id_success_count",
			Help:      "Total no of EventCountSummaryByCameraId success requests.",
		}),
		EventCountSummaryByEventTypeSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "event_count_summary_by_event_type_success_count",
			Help:      "Total no of EventCountSummaryByEventType success requests.",
		}),
		SummaryByCameraIdSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "summary_by_camera_id_success_count",
			Help:      "Total no of SummaryByCameraId success requests.",
		}),
		SummaryByEventTypeSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "summary_by_event_type_success_count",
			Help:      "Total no of SummaryByEventType success requests.",
		}),
	}
	promCamRegistry.MustRegister(
		m.EventCount,
		m.Info,
		m.DurationCountByEventType,
		m.DurationCountByCameraId,
		m.SuccessRequest,
		m.ErrorRequest,
		m.StoreEventDataSuccess,
		m.TotalEventCountByTypeSuccess,
		m.TotalEventCountByCameraIdSuccess,
		m.EventCountSummaryByEventTypeSuccess,
		m.EventCountSummaryByCameraIdSuccess,
		m.SummaryByEventTypeSuccess,
		m.SummaryByCameraIdSuccess,
	)
	return m
}
