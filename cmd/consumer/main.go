package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cwen96/social-analytics-pipeline/api"
	"github.com/cwen96/social-analytics-pipeline/internal/aggregator"
	"github.com/cwen96/social-analytics-pipeline/internal/kafka"
	"github.com/cwen96/social-analytics-pipeline/internal/metrics"
	"github.com/cwen96/social-analytics-pipeline/internal/store"
)

func main() {
	broker := envOrDefault("KAFKA_BROKER", "localhost:9092")
	topic := envOrDefault("KAFKA_TOPIC", "engagement-events")
	groupID := envOrDefault("KAFKA_GROUP_ID", "analytics-consumer")
	apiAddr := envOrDefault("API_ADDR", ":8080")
	windowSize := 1 * time.Minute

	// Initialize dependencies.
	metricsStore := store.New()
	promMetrics := metrics.New()
	agg := aggregator.New(metricsStore, windowSize, promMetrics)

	// Kafka consumer.
	consumer := kafka.NewConsumer(broker, topic, groupID, agg.Handle)

	// HTTP API server.
	handler := api.NewHandler(metricsStore)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	server := &http.Server{Addr: apiAddr, Handler: mux}

	// Graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API server listening on %s", apiAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server error: %v", err)
		}
	}()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("consumer error: %v", err)
		}
	}()

	<-sigCh
	log.Println("shutting down...")
	cancel()
	consumer.Close()
	server.Shutdown(context.Background())
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
