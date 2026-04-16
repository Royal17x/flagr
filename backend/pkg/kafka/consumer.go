package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type MessageHandler func(ctx context.Context, msg AuditMessage) error

type Consumer struct {
	reader     *kafka.Reader
	dlqWriter  *kafka.Writer
	maxRetries int
}

func NewConsumer(broker, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   AuditTopic,
		GroupID: groupID,

		StartOffset: kafka.LastOffset,

		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  time.Second,

		CommitInterval: time.Second,
	})

	dlqWriter := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        AuditDLQTopic,
		RequiredAcks: kafka.RequireOne,
	}

	return &Consumer{
		reader:     reader,
		dlqWriter:  dlqWriter,
		maxRetries: 3,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("kafka: failed to fetch message", "error", err)
			continue
		}

		if err := c.processWithRetry(ctx, msg, handler); err != nil {
			c.sendToDLQ(ctx, msg, err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			slog.Error("kafka: failed to commit offset", "error", err)
		}
	}
}

func (c *Consumer) processWithRetry(ctx context.Context, msg kafka.Message, handler MessageHandler) error {
	var auditMsg AuditMessage
	if err := json.Unmarshal(msg.Value, &auditMsg); err != nil {
		return fmt.Errorf("unmarshall: %w", err)
	}

	var lastErr error
	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		if err := handler(ctx, auditMsg); err != nil {
			lastErr = err
			slog.Warn("kafka: handler failed, retrying",
				"attempt", attempt,
				"max", c.maxRetries,
				"error", err,
				"message_id", auditMsg.ID,
			)

			time.Sleep(time.Duration(attempt*attempt) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	return lastErr
}

func (c *Consumer) sendToDLQ(ctx context.Context, msg kafka.Message, err error) {
	slog.Error("kafka: sending message to DLQ",
		"error", err,
		"topic", msg.Topic,
		"partition", msg.Partition,
		"offset", msg.Offset,
	)

	headers := append(msg.Headers,
		kafka.Header{Key: "dlq_reason", Value: []byte(err.Error())},
		kafka.Header{Key: "dlq_timestamp", Value: []byte(time.Now().UTC().Format(time.RFC3339))},
	)

	err = c.dlqWriter.WriteMessages(ctx, kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
	})
	if err != nil {
		slog.Error("kafka: failed to send to DLQ", "error", err)
	}
}

func (c *Consumer) Close() error {
	c.dlqWriter.Close()
	return c.reader.Close()
}
