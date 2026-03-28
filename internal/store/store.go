package store

import (
	"sync"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
)

// MetricsStore is an in-memory store for aggregated engagement metrics.
type MetricsStore struct {
	mu      sync.RWMutex
	buckets map[string]*models.AggregatedMetrics // key: "postID:platform:windowStart"
}

// BucketCount returns the number of active aggregation buckets.
func (s *MetricsStore) BucketCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buckets)
}

// New creates a new MetricsStore.
func New() *MetricsStore {
	return &MetricsStore{
		buckets: make(map[string]*models.AggregatedMetrics),
	}
}

func bucketKey(postID, platform string, windowStart time.Time) string {
	return postID + ":" + platform + ":" + windowStart.Format(time.RFC3339)
}

// Record adds an engagement event to the appropriate time-window bucket.
func (s *MetricsStore) Record(event models.EngagementEvent, windowSize time.Duration) {
	windowStart := event.Timestamp.Truncate(windowSize)
	windowEnd := windowStart.Add(windowSize)
	key := bucketKey(event.PostID, event.Platform, windowStart)

	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, exists := s.buckets[key]
	if !exists {
		bucket = &models.AggregatedMetrics{
			PostID:      event.PostID,
			Platform:    event.Platform,
			WindowStart: windowStart,
			WindowEnd:   windowEnd,
			Counts:      make(map[models.EventType]int),
		}
		s.buckets[key] = bucket
	}

	bucket.Counts[event.Type]++
	bucket.Total++
}

// Query returns aggregated metrics matching the given filters.
func (s *MetricsStore) Query(q models.MetricsQuery) []models.AggregatedMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.AggregatedMetrics
	for _, bucket := range s.buckets {
		if q.PostID != "" && bucket.PostID != q.PostID {
			continue
		}
		if q.Platform != "" && bucket.Platform != q.Platform {
			continue
		}
		if !q.From.IsZero() && bucket.WindowEnd.Before(q.From) {
			continue
		}
		if !q.To.IsZero() && bucket.WindowStart.After(q.To) {
			continue
		}
		results = append(results, *bucket)
	}
	return results
}

// All returns every bucket in the store.
func (s *MetricsStore) All() []models.AggregatedMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]models.AggregatedMetrics, 0, len(s.buckets))
	for _, bucket := range s.buckets {
		results = append(results, *bucket)
	}
	return results
}
