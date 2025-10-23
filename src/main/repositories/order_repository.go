package repositories

import (
	"context"

	domain "order-management-ms/src/main/models/datastore"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Create saves a new order to the database
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	// FindByID finds an order by its order ID
	FindByID(ctx context.Context, orderID string) (*domain.Order, error)

	// UpdateStatus updates the status of an order
	UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error

	// List returns a list of orders with pagination and filtering
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*domain.Order, error)
}
