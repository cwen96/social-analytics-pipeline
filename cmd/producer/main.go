package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/kafka"
)

func main() {
	broker := envOrDefault("KAFKA_BROKER", "localhost:9092")
	topic := envOrDefault("KAFKA_TOPIC", "engagement-events")
	interval := 500 * time.Millisecond

	producer := kafka.NewProducer(broker, topic)
	defer producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("shutting down producer...")
		cancel()
	}()

	log.Println("starting event producer...")
	if err := producer.Start(ctx, interval); err != nil {
		log.Fatalf("producer error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
