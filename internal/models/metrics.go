package models

import "time"

// AggregatedMetrics holds engagement counts for a given time window.
type AggregatedMetrics struct {
	PostID    string           `json:"post_id"`
	Platform  string           `json:"platform"`
	WindowStart time.Time      `json:"window_start"`
	WindowEnd   time.Time      `json:"window_end"`
	Counts    map[EventType]int `json:"counts"`
	Total     int              `json:"total"`
}

// MetricsQuery represents filters for querying aggregated metrics.
type MetricsQuery struct {
	PostID   string `json:"post_id"`
	Platform string `json:"platform"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
}
