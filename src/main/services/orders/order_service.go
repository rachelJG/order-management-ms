package orders

import (
	"context"
	"encoding/json"
	models "order-management-ms/src/main/models/api"
	domain "order-management-ms/src/main/models/datastore"
	kafkaDto "order-management-ms/src/main/models/kafka"
	errors "order-management-ms/src/main/pkg/customerrors"
	"time"

	"go.uber.org/zap"
)

// Service defines the interface for order operations
type Service interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*models.OrderResponse, error)
	GetOrder(ctx context.Context, orderID string) (*models.OrderResponse, error)
	ListOrders(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus domain.OrderStatus) error
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*models.OrderResponse, error) {
	// Set default values
	order.Status = domain.StatusNew
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Save to database
	newOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err), zap.String("order_id", order.ID.Hex()))
		return nil, err
	}

	return models.NewOrderResponse(newOrder), nil
}

// GetOrder retrieves an order by ID with caching
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	// Try to get from cache first
	cachedOrder, err := s.getFromCache(ctx, orderID)
	if err == nil && cachedOrder != nil {
		return cachedOrder, nil
	}

	// If not in cache, get from database
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to find order",
			zap.Error(err),
			zap.String("order_id", orderID),
		)
		return nil, errors.ErrOrderNotFound
	}

	return models.NewOrderResponse(order), nil
}

// ListOrders retrieves a list of orders with optional filters
func (s *OrderService) ListOrders(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.OrderResponse, error) {
	orders, err := s.repo.List(ctx, filters, page, limit)
	if err != nil {
		s.logger.Error("Failed to list orders", zap.Error(err))
		return nil, err
	}

	var orderResponses []*models.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, models.NewOrderResponse(order))
	}

	return orderResponses, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, newStatus domain.OrderStatus) error {
	// Get current order
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to find order",
			zap.Error(err),
			zap.String("order_id", orderID),
		)
		return errors.ErrOrderNotFound
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, newStatus) {
		s.logger.Warn("Invalid status transition",
			zap.String("order_id", orderID),
			zap.String("current_status", string(order.Status)),
			zap.String("new_status", string(newStatus)),
		)
		return errors.ErrInvalidTransition
	}

	// Save old status for event
	oldStatus := order.Status

	// Update status
	order.Status = newStatus
	order.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		s.logger.Error("Failed to update order status",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.String("new_status", string(newStatus)),
		)
		return errors.ErrInternalServer
	}

	event := kafkaDto.NewOrderStatusChangedEvent(orderID, oldStatus, newStatus)

	if err := s.eventPublisher.PublishOrderStatusChanged(ctx, event); err != nil {
		s.logger.Error("Failed to publish order status changed event",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.Any("event", event),
		)
	}

	// Invalidate cache for this order
	cacheKey := "order:" + orderID
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		s.logger.Error("Failed to invalidate cache for order",
			zap.Error(err),
			zap.String("order_id", orderID),
		)
	}

	return nil
}

// SaveOrderInCache saves an order in the cache with a TTL of 60 seconds
func (s *OrderService) SaveOrderInCache(ctx context.Context, order *domain.Order) error {
	if order == nil {
		return nil
	}

	cacheKey := "order:" + order.OrderID
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, cacheKey, string(orderJSON), 60*time.Second)
}

// getFromCache attempts to retrieve an order from the cache
func (s *OrderService) getFromCache(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	cacheKey := "order:" + orderID
	val, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, err
	}

	return models.NewOrderResponse(&order), nil
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
