package kafka

import (
	"context"
	"order-management-ms/src/main/models/datastore"
	"time"
)

// EventPublisher defines the contract for publishing events in the system
type EventPublisher interface {
	// PublishOrderStatusChanged publishes an order status changed event
	PublishOrderStatusChanged(ctx context.Context, event OrderStatusChangedEvent) error
}

type OrderStatusChangedEvent struct {
	OrderID   string                `json:"order_id"`
	OldStatus datastore.OrderStatus `json:"old_status"`
	NewStatus datastore.OrderStatus `json:"new_status"`
	Timestamp string                `json:"timestamp"`
}

func NewOrderStatusChangedEvent(orderID string, oldStatus, newStatus datastore.OrderStatus) OrderStatusChangedEvent {
	return OrderStatusChangedEvent{
		OrderID:   orderID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
