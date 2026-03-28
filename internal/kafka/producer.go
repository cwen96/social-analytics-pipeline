package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cwen96/social-analytics-pipeline/internal/models"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

// Producer simulates social media engagement events and writes them to Kafka.
type Producer struct {
	writer *kafkago.Writer
}

// NewProducer creates a Kafka producer for the given topic and broker.
func NewProducer(broker, topic string) *Producer {
	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(broker),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}

	return &Producer{writer: writer}
}

var (
	platforms  = []string{"twitter", "instagram", "facebook", "linkedin", "tiktok"}
	eventTypes = []models.EventType{
		models.EventLike,
		models.EventShare,
		models.EventClick,
		models.EventComment,
		models.EventRepost,
	}
)

// GenerateEvent creates a random engagement event.
func GenerateEvent() models.EngagementEvent {
	return models.EngagementEvent{
		ID:        uuid.New().String(),
		PostID:    fmt.Sprintf("post-%d", rand.Intn(50)+1),
		UserID:    fmt.Sprintf("user-%d", rand.Intn(1000)+1),
		Platform:  platforms[rand.Intn(len(platforms))],
		Type:      eventTypes[rand.Intn(len(eventTypes))],
		Timestamp: time.Now().UTC(),
	}
}

// Start begins producing simulated events at the given interval.
func (p *Producer) Start(ctx context.Context, interval time.Duration) error {
	log.Printf("kafka producer started, emitting every %s", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			event := GenerateEvent()
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("error marshalling event: %v", err)
				continue
			}

			err = p.writer.WriteMessages(ctx, kafkago.Message{
				Key:   []byte(event.PostID),
				Value: data,
			})
			if err != nil {
				log.Printf("error writing message: %v", err)
				continue
			}

			log.Printf("produced event: %s %s on %s for %s", event.Type, event.ID[:8], event.Platform, event.PostID)
		}
	}
}

// Close shuts down the producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
