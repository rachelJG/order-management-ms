package kafka

import (
	"context"
	"encoding/json"
	kafkaDto "order-management-ms/src/main/models/kafka"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Publisher struct {
	writer *kafka.Writer
	logger *zap.Logger
	topic  string
}

func NewPublisher(brokers []string, topic string, logger *zap.Logger) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
		topic:  topic,
		logger: logger,
	}
}

func (k *Publisher) PublishOrderStatusChanged(ctx context.Context, event kafkaDto.OrderStatusChangedEvent) error {
	// Serialize the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		k.logger.Error("Failed to serialize event",
			zap.Error(err),
			zap.String("order_id", event.OrderID),
		)
		return err
	}

	// Publish the message to Kafka
	err = k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: eventJSON,
	})

	if err != nil {
		k.logger.Error("Failed to publish message to Kafka",
			zap.Error(err),
			zap.String("topic", k.topic),
			zap.String("order_id", event.OrderID),
		)
		return err
	}

	k.logger.Debug("Event published successfully",
		zap.String("topic", k.topic),
		zap.String("order_id", event.OrderID),
		zap.String("old_status", string(event.OldStatus)),
		zap.String("new_status", string(event.NewStatus)),
	)

	return nil
}

// Close closes the connection to Kafka
func (k *Publisher) Close() error {
	if k.writer != nil {
		return k.writer.Close()
	}
	return nil
}
