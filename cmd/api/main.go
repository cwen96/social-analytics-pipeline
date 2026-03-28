package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cwen96/social-analytics-pipeline/api"
	"github.com/cwen96/social-analytics-pipeline/internal/store"
)

// Standalone API server (useful for development without Kafka).
func main() {
	addr := envOrDefault("API_ADDR", ":8080")

	metricsStore := store.New()
	handler := api.NewHandler(metricsStore)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Printf("standalone API server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
