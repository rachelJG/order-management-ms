package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Producer struct {
	writer Writer
	logger *zap.Logger
}

func NewProducer(brokers []string, topic string, logger *zap.Logger) *Producer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	return &Producer{
		writer: w,
		logger: logger,
	}
}

func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{Key: key, Value: value}
	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.logger.Error("failed to publish message", zap.Error(err))
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
