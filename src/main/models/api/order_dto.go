package api

import (
	"time"
)

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Items      []Items `json:"items" binding:"required,min=1"`
}

type Items struct {
	Sku      string  `json:"sku" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Price    float64 `json:"price,omitempty"`
}

// ListOrdersRequest represents the query parameters for listing orders
type ListOrdersRequest struct {
	Status     string `form:"status"`
	CustomerID string `form:"customer_id"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}

// OrderResponse represents the response body for an order
type OrderResponse struct {
	CustomerID string    `json:"customer_id"`
	OrderID    string    `json:"order_id"`
	Status     string    `json:"status"`
	Items      []Items   `json:"items"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UpdateOrderStatusRequest represents the request body for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
