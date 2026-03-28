package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/cwen96/social-analytics-pipeline/internal/store"
)

func setupTestServer() (*Handler, *http.ServeMux) {
	s := store.New()
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	s.Record(models.EngagementEvent{ID: "1", PostID: "post-1", Platform: "twitter", Type: models.EventLike, Timestamp: now}, time.Minute)
	s.Record(models.EngagementEvent{ID: "2", PostID: "post-1", Platform: "twitter", Type: models.EventShare, Timestamp: now}, time.Minute)
	s.Record(models.EngagementEvent{ID: "3", PostID: "post-2", Platform: "instagram", Type: models.EventClick, Timestamp: now}, time.Minute)

	h := NewHandler(s)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestHealthEndpoint(t *testing.T) {
	_, mux := setupTestServer()
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %s", resp["status"])
	}
}

func TestAllMetricsEndpoint(t *testing.T) {
	_, mux := setupTestServer()
	req := httptest.NewRequest("GET", "/metrics/all", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var results []models.AggregatedMetrics
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(results))
	}
}

func TestQueryByPostID(t *testing.T) {
	_, mux := setupTestServer()
	req := httptest.NewRequest("GET", "/metrics/query?post_id=post-1", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var results []models.AggregatedMetrics
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestQueryByPlatform(t *testing.T) {
	_, mux := setupTestServer()
	req := httptest.NewRequest("GET", "/metrics/query?platform=instagram", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	var results []models.AggregatedMetrics
	json.NewDecoder(w.Body).Decode(&results)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestQueryInvalidFromParam(t *testing.T) {
	_, mux := setupTestServer()
	req := httptest.NewRequest("GET", "/metrics/query?from=not-a-date", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
