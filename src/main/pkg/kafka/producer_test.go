package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockKafkaWriter struct {
	mock.Mock
}

func (m *mockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *mockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewProducer(t *testing.T) {
	logger := zap.NewNop()
	brokers := []string{"localhost:9092"}
	topic := "test-topic"

	producer := NewProducer(brokers, topic, logger)

	assert.NotNil(t, producer)
	assert.NotNil(t, producer.writer)
	assert.Equal(t, logger, producer.logger)
}

func TestPublish_Success(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()

	producer := &Producer{
		writer: mockWriter,
		logger: logger,
	}

	ctx := context.Background()
	key := []byte("test-key")
	value := []byte("test-value")

	mockWriter.On("WriteMessages", ctx, mock.AnythingOfType("[]kafka.Message")).Return(nil)

	err := producer.Publish(ctx, key, value)

	assert.NoError(t, err)
	mockWriter.AssertExpectations(t)
}

func TestPublish_Error(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()
	producer := &Producer{
		writer: mockWriter,
		logger: logger,
	}

	ctx := context.Background()
	expectedErr := errors.New("kafka write error")

	mockWriter.On("WriteMessages", ctx, mock.Anything).Return(expectedErr)

	err := producer.Publish(ctx, []byte("key"), []byte("value"))

	assert.ErrorIs(t, err, expectedErr)
	mockWriter.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()
	producer := &Producer{
		writer: mockWriter,
		logger: logger,
	}

	t.Run("success", func(t *testing.T) {
		mockWriter.On("Close").Return(nil)
		err := producer.Close()
		assert.NoError(t, err)
		mockWriter.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("close error")
		mockWriter.ExpectedCalls = nil
		mockWriter.On("Close").Return(expectedErr)
		err := producer.Close()
		assert.ErrorIs(t, err, expectedErr)
		mockWriter.AssertExpectations(t)
	})
}
