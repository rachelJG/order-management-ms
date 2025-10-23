package requests

import "order-management-ms/src/main/domain"

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id" binding:"required"`
	Items      []domain.OrderItem `json:"items" binding:"required,min=1"`
}

// ListOrdersRequest represents the query parameters for listing orders
type ListOrdersRequest struct {
	Status     *domain.OrderStatus `form:"status"`
	CustomerID string              `form:"customer_id"`
	Page       int                 `form:"page,default=1"`
	Limit      int                 `form:"limit,default=10"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" binding:"required,oneof=processing shipped delivered cancelled"`
}
