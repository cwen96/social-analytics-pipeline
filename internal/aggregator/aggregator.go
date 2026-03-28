package aggregator

import (
	"log"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/metrics"
	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/cwen96/social-analytics-pipeline/internal/store"
)

// Aggregator processes engagement events and stores them in time-window buckets.
type Aggregator struct {
	store      *store.MetricsStore
	windowSize time.Duration
	metrics    *metrics.Metrics
}

// New creates an Aggregator with the given window size.
func New(s *store.MetricsStore, windowSize time.Duration, m *metrics.Metrics) *Aggregator {
	return &Aggregator{
		store:      s,
		windowSize: windowSize,
		metrics:    m,
	}
}

// Handle processes a single engagement event.
func (a *Aggregator) Handle(event models.EngagementEvent) {
	start := time.Now()
	a.store.Record(event, a.windowSize)
	a.metrics.ProcessingTime.WithLabelValues(string(event.Type)).Observe(time.Since(start).Seconds())
	a.metrics.ActiveBuckets.Set(float64(a.store.BucketCount()))
	a.metrics.RecordEvent(event)
	log.Printf("aggregated event: %s %s on %s for %s", event.Type, event.ID[:8], event.Platform, event.PostID)
}
