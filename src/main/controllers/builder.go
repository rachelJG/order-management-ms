package controllers

import (
	"order-management-ms/src/main/services/orders"

	"go.uber.org/zap"
)

type OrderController struct {
	service *orders.OrderService
	logger  *zap.Logger
}

func NewOrderController(service *orders.OrderService, logger *zap.Logger) *OrderController {
	return &OrderController{
		service: service,
		logger:  logger,
	}
}
