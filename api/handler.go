package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/cwen96/social-analytics-pipeline/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler holds dependencies for the HTTP API.
type Handler struct {
	store *store.MetricsStore
}

// NewHandler creates an API handler backed by the given store.
func NewHandler(s *store.MetricsStore) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes sets up the HTTP routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /metrics/query", h.queryMetrics)
	mux.HandleFunc("GET /metrics/all", h.allMetrics)
	mux.HandleFunc("GET /health", h.health)
	mux.Handle("GET /prometheus", promhttp.Handler())
}

func (h *Handler) queryMetrics(w http.ResponseWriter, r *http.Request) {
	q := models.MetricsQuery{
		PostID:   r.URL.Query().Get("post_id"),
		Platform: r.URL.Query().Get("platform"),
	}

	if from := r.URL.Query().Get("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			http.Error(w, "invalid 'from' parameter: use RFC3339 format", http.StatusBadRequest)
			return
		}
		q.From = t
	}
	if to := r.URL.Query().Get("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			http.Error(w, "invalid 'to' parameter: use RFC3339 format", http.StatusBadRequest)
			return
		}
		q.To = t
	}

	results := h.store.Query(q)
	writeJSON(w, results)
}

func (h *Handler) allMetrics(w http.ResponseWriter, _ *http.Request) {
	results := h.store.All()
	writeJSON(w, results)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
