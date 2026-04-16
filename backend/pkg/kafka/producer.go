package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

type AuditMessage struct {
	ID         uuid.UUID `json:"id"`
	Version    int       `json:"version"`
	Action     string    `json:"action"`
	ActorID    uuid.UUID `json:"actor_id"`
	ResourceID uuid.UUID `json:"resource_id"`
	ProjectID  uuid.UUID `json:"project_id"`
	Payload    any       `json:"payload,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewProducer(broker string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: AuditTopic,

			Balancer: &kafka.Murmur2Balancer{},

			RequiredAcks: kafka.RequireOne,

			Async: false,

			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,

			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,

			AllowAutoTopicCreation: false,

			Logger:      kafka.LoggerFunc(func(msg string, args ...interface{}) {}),
			ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {}),
		},
	}
}

func (p *Producer) PublishAuditEvent(ctx context.Context, event AuditMessage) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("producer.PublishAuditEvent: marshall: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ResourceID.String()),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Action)},
			{Key: "version", Value: []byte("1")},
			{Key: "produced_at", Value: []byte(time.Now().UTC().Format(time.RFC3339))},
		},
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
