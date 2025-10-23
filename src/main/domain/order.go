package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusInProgress OrderStatus = "IN_PROGRESS"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID string             `bson:"customer_id" json:"customer_id"`
	Status     OrderStatus        `bson:"status" json:"status"`
	Items      []OrderItem        `bson:"items" json:"items"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Sku       string  `bson:"sku" json:"sku"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
}

type OrderEvent struct {
	EventType  string      `json:"event_type"`
	OrderID    string      `json:"order_id"`
	CustomerID string      `json:"customer_id"`
	OldStatus  OrderStatus `json:"old_status,omitempty"`
	NewStatus  OrderStatus `json:"new_status"`
	Items      []OrderItem `json:"items"`
	Timestamp  time.Time   `json:"timestamp"`
}

// OrderRepository is the interface for order repository, is used to abstract the persistence layer
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*Order, error)
	FindByStatus(ctx context.Context, status OrderStatus) ([]*Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
	Delete(ctx context.Context, id string) error
}

// OrderService is the interface for order service, is used to abstract the business logic layer
type OrderService interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id string) (*Order, error)
	GetOrdersByCustomer(ctx context.Context, customerID string) ([]*Order, error)
	GetOrdersByStatus(ctx context.Context, status string) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, id string, newStatus string) error
	CancelOrder(ctx context.Context, id string) error
}
