package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/segmentio/kafka-go"
)

// Consumer reads engagement events from a Kafka topic.
type Consumer struct {
	reader  *kafka.Reader
	handler func(models.EngagementEvent)
}

// NewConsumer creates a Kafka consumer for the given topic and broker.
func NewConsumer(broker, topic, groupID string, handler func(models.EngagementEvent)) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

// Start begins consuming messages. Blocks until the context is cancelled.
func (c *Consumer) Start(ctx context.Context) error {
	log.Println("kafka consumer started")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // graceful shutdown
			}
			log.Printf("error reading message: %v", err)
			continue
		}

		var event models.EngagementEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("error unmarshalling event: %v", err)
			continue
		}

		c.handler(event)
	}
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
