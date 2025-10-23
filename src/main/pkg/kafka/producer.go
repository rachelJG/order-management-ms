package kafka

import (
	"context"
	"encoding/json"

	kafkaDto "order-management-ms/src/main/models/kafka"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Producer implements the events.EventPublisher interface
// and provides methods to publish messages to Kafka
type Producer struct {
	writer Writer
	logger *zap.Logger
	topic  string
}

// Ensure Producer implements domain.EventPublisher
var _ kafkaDto.EventPublisher = (*Producer)(nil)

// NewProducer creates a new Kafka producer instance
func NewProducer(brokers []string, topic string, logger *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	return &Producer{
		writer: w,
		logger: logger,
		topic:  topic,
	}
}

// PublishOrderStatusChanged implements the EventPublisher interface
func (p *Producer) PublishOrderStatusChanged(ctx context.Context, event kafkaDto.OrderStatusChangedEvent) error {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Error serializing event",
			zap.Error(err),
			zap.String("order_id", event.OrderID),
		)
		return err
	}

	// Publish message to Kafka
	msg := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: eventJSON,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.logger.Error("Failed to publish message to Kafka",
			zap.Error(err),
			zap.String("order_id", event.OrderID),
			zap.String("topic", p.topic),
		)
		return err
	}

	p.logger.Debug("Successfully published order status change event",
		zap.String("order_id", event.OrderID),
		zap.String("old_status", string(event.OldStatus)),
		zap.String("new_status", string(event.NewStatus)),
	)

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
