package store

import (
	"testing"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
)

func TestRecord_CreatesNewBucket(t *testing.T) {
	s := New()
	now := time.Date(2025, 1, 1, 12, 0, 30, 0, time.UTC)
	event := models.EngagementEvent{
		ID:        "test-1",
		PostID:    "post-1",
		UserID:    "user-1",
		Platform:  "twitter",
		Type:      models.EventLike,
		Timestamp: now,
	}

	s.Record(event, 1*time.Minute)

	results := s.All()
	if len(results) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(results))
	}
	if results[0].Total != 1 {
		t.Errorf("expected total 1, got %d", results[0].Total)
	}
	if results[0].Counts[models.EventLike] != 1 {
		t.Errorf("expected 1 like, got %d", results[0].Counts[models.EventLike])
	}
	// Window should be truncated to the minute.
	expectedStart := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	if !results[0].WindowStart.Equal(expectedStart) {
		t.Errorf("expected window start %v, got %v", expectedStart, results[0].WindowStart)
	}
}

func TestRecord_AggregatesSameBucket(t *testing.T) {
	s := New()
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	for i := 0; i < 5; i++ {
		s.Record(models.EngagementEvent{
			ID:        "test",
			PostID:    "post-1",
			Platform:  "instagram",
			Type:      models.EventClick,
			Timestamp: base.Add(time.Duration(i) * time.Second),
		}, 1*time.Minute)
	}

	results := s.All()
	if len(results) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(results))
	}
	if results[0].Total != 5 {
		t.Errorf("expected total 5, got %d", results[0].Total)
	}
}

func TestRecord_SeparatesBucketsByWindow(t *testing.T) {
	s := New()
	s.Record(models.EngagementEvent{
		ID: "a", PostID: "post-1", Platform: "twitter", Type: models.EventLike,
		Timestamp: time.Date(2025, 1, 1, 12, 0, 30, 0, time.UTC),
	}, 1*time.Minute)
	s.Record(models.EngagementEvent{
		ID: "b", PostID: "post-1", Platform: "twitter", Type: models.EventLike,
		Timestamp: time.Date(2025, 1, 1, 12, 1, 30, 0, time.UTC),
	}, 1*time.Minute)

	results := s.All()
	if len(results) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(results))
	}
}

func TestQuery_FiltersByPostID(t *testing.T) {
	s := New()
	now := time.Now().UTC()
	s.Record(models.EngagementEvent{ID: "a", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: now}, time.Minute)
	s.Record(models.EngagementEvent{ID: "b", PostID: "post-2", Platform: "twitter", Type: models.EventLike, Timestamp: now}, time.Minute)

	results := s.Query(models.MetricsQuery{PostID: "post-1"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].PostID != "post-1" {
		t.Errorf("expected post-1, got %s", results[0].PostID)
	}
}

func TestQuery_FiltersByPlatform(t *testing.T) {
	s := New()
	now := time.Now().UTC()
	s.Record(models.EngagementEvent{ID: "a", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: now}, time.Minute)
	s.Record(models.EngagementEvent{ID: "b", PostID: "post-1", Platform: "instagram", Type: models.EventLike, Timestamp: now}, time.Minute)

	results := s.Query(models.MetricsQuery{Platform: "instagram"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Platform != "instagram" {
		t.Errorf("expected instagram, got %s", results[0].Platform)
	}
}

func TestQuery_FiltersByTimeRange(t *testing.T) {
	s := New()
	t1 := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC)
	t3 := time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC)

	s.Record(models.EngagementEvent{ID: "a", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: t1}, time.Minute)
	s.Record(models.EngagementEvent{ID: "b", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: t2}, time.Minute)
	s.Record(models.EngagementEvent{ID: "c", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: t3}, time.Minute)

	results := s.Query(models.MetricsQuery{
		From: time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC),
		To:   time.Date(2025, 1, 1, 13, 30, 0, 0, time.UTC),
	})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}
