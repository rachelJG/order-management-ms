package orders

import (
	"context"
	"order-management-ms/src/main/domain"
	errors "order-management-ms/src/main/pkg/customerrors"
	"time"

	"go.uber.org/zap"
)

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {

	// Validate order
	if err := s.validateOrder(order); err != nil {
		s.logger.Error("Invalid order", zap.Error(err), zap.String("order_id", order.ID.Hex()))
		return errors.ErrInvalidOrder
	}

	// Set default values
	order.Status = domain.StatusNew
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Create(ctx, order); err != nil {
		s.logger.Error("Failed to create order", zap.Error(err), zap.String("order_id", order.ID.Hex()))
		return err
	}

	return nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get order", zap.Error(err), zap.String("order_id", id))
		return nil, errors.ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context, filters map[string]interface{}) ([]domain.Order, error) {
	// Validate status if provided
	if status, ok := filters["status"].(string); ok {
		if !isValidStatus(domain.OrderStatus(status)) {
			return nil, errors.ErrInvalidStatus
		}
	}

	orders, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func isValidStatus(status domain.OrderStatus) bool {
	switch status {
	case domain.StatusNew, domain.StatusInProgress, domain.StatusDelivered, domain.StatusCancelled:
		return true
	default:
		return false
	}
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, newStatus domain.OrderStatus) error {
	// Get current order
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find order", zap.Error(err), zap.String("order_id", id))
		return errors.ErrOrderNotFound
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, newStatus) {
		s.logger.Warn("Invalid status transition",
			zap.String("order_id", id),
			zap.String("current_status", string(order.Status)),
			zap.String("new_status", string(newStatus)),
		)
		return errors.ErrInvalidTransition
	}

	// Update status
	order.Status = newStatus
	order.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		s.logger.Error("Failed to update order status",
			zap.Error(err),
			zap.String("order_id", id),
			zap.String("status", string(newStatus)),
		)
		return err
	}

	return nil
}

// validateOrder validates the order
func (s *OrderService) validateOrder(order *domain.Order) error {
	if order == nil {
		return errors.ErrInvalidOrder
	}

	if len(order.Items) == 0 {
		return errors.ErrItemsIsRequired
	}

	return nil
}

// isValidStatusTransition checks if the status transition is valid
func isValidStatusTransition(current, newStatus domain.OrderStatus) bool {
	switch current {
	case domain.StatusNew:
		return newStatus == domain.StatusInProgress || newStatus == domain.StatusCancelled
	case domain.StatusInProgress:
		return newStatus == domain.StatusDelivered || newStatus == domain.StatusCancelled
	case domain.StatusDelivered, domain.StatusCancelled:
		return false
	default:
		return false
	}
}
