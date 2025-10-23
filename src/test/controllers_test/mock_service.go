package controllers_test

import (
	"context"

	"order-management-ms/src/main/models/api"
	dm "order-management-ms/src/main/models/datastore"

	"github.com/stretchr/testify/mock"
)

// mockOrderService is a mock implementation of the OrderService
type mockOrderService struct {
	mock.Mock
}

// CreateOrder mocks the CreateOrder method
func (m *mockOrderService) CreateOrder(ctx context.Context, order *dm.Order) (*api.OrderResponse, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.OrderResponse), args.Error(1)
}

// GetOrder mocks the GetOrder method
func (m *mockOrderService) GetOrder(ctx context.Context, orderID string) (*api.OrderResponse, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.OrderResponse), args.Error(1)
}

// ListOrders mocks the ListOrders method
func (m *mockOrderService) ListOrders(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*api.OrderResponse, error) {
	args := m.Called(ctx, filters, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*api.OrderResponse), args.Error(1)
}

// UpdateOrderStatus mocks the UpdateOrderStatus method
func (m *mockOrderService) UpdateOrderStatus(ctx context.Context, orderID string, status dm.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}
