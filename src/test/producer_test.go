package test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"order-management-ms/src/main/models/datastore"
	"order-management-ms/src/main/pkg/kafka"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// mockKafkaWriter simulates the kafka.Writer interface for unit testing.
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

// --- Tests ---

func TestNewProducer(t *testing.T) {
	logger := zap.NewNop()
	mockWriter := new(mockKafkaWriter)
	topic := "test-topic"

	producer := kafka.NewProducer(mockWriter, logger, topic)

	assert.NotNil(t, producer)
	assert.Equal(t, logger, producer.Logger)
	assert.Equal(t, topic, producer.Topic)
}

func TestPublishOrderStatusChanged_Success(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()
	producer := kafka.NewProducer(mockWriter, logger, "test-topic")

	ctx := context.Background()
	event := datastore.OrderStatusChangedEvent{
		OrderID:   "test-order-123",
		OldStatus: datastore.StatusNew,
		NewStatus: datastore.StatusInProgress,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	expectedValue, _ := json.Marshal(event)
	expectedMsg := kafka.Message{
		Key:   []byte(event.OrderID),
		Value: expectedValue,
	}

	// ✅ Use mock.MatchedBy to compare slice contents instead of slice references.
	mockWriter.On("WriteMessages", ctx, mock.MatchedBy(func(msgs []kafka.Message) bool {
		return len(msgs) == 1 &&
			string(msgs[0].Key) == event.OrderID &&
			string(msgs[0].Value) == string(expectedValue)
	})).Return(nil)

	err := producer.PublishOrderStatusChanged(ctx, event)

	assert.NoError(t, err)
	mockWriter.AssertExpectations(t)
}

func TestPublishOrderStatusChanged_Error(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()
	producer := kafka.NewProducer(mockWriter, logger, "test-topic")

	ctx := context.Background()
	event := datastore.OrderStatusChangedEvent{
		OrderID:   "test-order-123",
		OldStatus: datastore.StatusNew,
		NewStatus: datastore.StatusInProgress,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	expectedErr := errors.New("kafka write error")

	expectedValue, _ := json.Marshal(event)

	// ✅ Using mock.MatchedBy avoids direct slice comparison.
	mockWriter.On("WriteMessages", ctx, mock.MatchedBy(func(msgs []kafka.Message) bool {
		return len(msgs) == 1 &&
			string(msgs[0].Key) == event.OrderID &&
			string(msgs[0].Value) == string(expectedValue)
	})).Return(expectedErr)

	err := producer.PublishOrderStatusChanged(ctx, event)

	assert.ErrorIs(t, err, expectedErr)
	mockWriter.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	mockWriter := new(mockKafkaWriter)
	logger := zap.NewNop()
	producer := &kafka.Producer{
		Writer: mockWriter,
		Logger: logger,
	}

	t.Run("success", func(t *testing.T) {
		// Clear previous expectations before each subtest
		mockWriter.ExpectedCalls = nil
		mockWriter.On("Close").Return(nil)

		err := producer.Close()
		assert.NoError(t, err)
		mockWriter.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		// Reset mock to avoid test interference
		mockWriter.ExpectedCalls = nil
		expectedErr := errors.New("close error")
		mockWriter.On("Close").Return(expectedErr)

		err := producer.Close()
		// ✅ Use assert.ErrorIs for idiomatic Go error comparison
		assert.ErrorIs(t, err, expectedErr)
		mockWriter.AssertExpectations(t)
	})
}
