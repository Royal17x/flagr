package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

const (
	AuditTopic    = "flag.audit"
	AuditDLQTopic = "flag.audit.dlq"
)

func EnsureTopics(ctx context.Context, broker string, replicationFactor int16) error {
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		return fmt.Errorf("kafka.EnsureTopics: dial: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("kafka.EnsureTopics: controller: %w", err)
	}

	controllerConn, err := kafka.DialContext(ctx, "tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		return fmt.Errorf("kafka.EnsureTopics: controller dial: %w", err)
	}
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{
			Topic:             AuditTopic,
			NumPartitions:     3,
			ReplicationFactor: int(replicationFactor),
		},
		{
			Topic:             AuditDLQTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		var kafkaErr kafka.Error
		if errors.As(err, &kafkaErr) && errors.Is(kafkaErr, kafka.TopicAlreadyExists) {
			return nil
		}
		return fmt.Errorf("kafka.EnsureTopics: create topics: %w", err)
	}
	return nil
}
