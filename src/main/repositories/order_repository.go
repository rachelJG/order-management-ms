package repositories

import (
	"context"

	"order-management-ms/src/main/domain"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	// Create saves a new order to the database
	Create(ctx context.Context, order *domain.Order) error

	// FindByID finds an order by its ID
	FindByID(ctx context.Context, id string) (*domain.Order, error)

	// UpdateStatus updates the status of an order
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error

	// List returns a list of orders with pagination and filtering
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*domain.Order, error)
}
