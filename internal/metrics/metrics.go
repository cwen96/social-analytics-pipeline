package metrics

import (
	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds Prometheus metrics for the analytics pipeline.
type Metrics struct {
	EventsConsumed *prometheus.CounterVec
	EventsProduced prometheus.Counter
	ProcessingTime *prometheus.HistogramVec
	ActiveBuckets  prometheus.Gauge
}

// New registers and returns all Prometheus metrics.
func New() *Metrics {
	return &Metrics{
		EventsConsumed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "analytics_events_consumed_total",
			Help: "Total number of engagement events consumed from Kafka",
		}, []string{"platform", "event_type"}),

		EventsProduced: promauto.NewCounter(prometheus.CounterOpts{
			Name: "analytics_events_produced_total",
			Help: "Total number of simulated events produced to Kafka",
		}),

		ProcessingTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "analytics_event_processing_seconds",
			Help:    "Time taken to process an engagement event",
			Buckets: prometheus.DefBuckets,
		}, []string{"event_type"}),

		ActiveBuckets: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "analytics_active_buckets",
			Help: "Number of active aggregation time-window buckets",
		}),
	}
}

// RecordEvent increments the consumed counter for the event's platform and type.
func (m *Metrics) RecordEvent(event models.EngagementEvent) {
	m.EventsConsumed.WithLabelValues(event.Platform, string(event.Type)).Inc()
}
