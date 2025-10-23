package controllers

import (
	"net/http"
	"order-management-ms/src/main/domain"
	"order-management-ms/src/main/models/dto/requests"
	errors "order-management-ms/src/main/pkg/customerrors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrder handles the creation of a new order
// @Summary Create a new order
// @Description Creates a new order with the provided items
// @Tags orders
// @Accept json
// @Produce json
// @Param input body CreateOrderRequest true "Order data"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req requests.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	order := &domain.Order{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	}

	if err := c.service.CreateOrder(ctx.Request.Context(), order); err != nil {
		c.logger.Error("Failed to create order", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrFailedToCreateOrder.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

// GetOrder handles retrieving an order by ID
// @Summary Get an order by ID
// @Description Retrieves details of a specific order
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders/{id} [get]
func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	order, err := c.service.GetOrder(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to get order", zap.Error(err), zap.String("order_id", id))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrOrderNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// ListOrders handles listing orders
// @Summary List orders
// @Description Lists all orders
// @Tags orders
// @Produce json
// @Success 200 {object} []domain.Order
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders [get]
func (c *OrderController) ListOrders(ctx *gin.Context) {
	// Get query parameters
	status := ctx.Query("status")
	customerID := ctx.Query("customer_id")

	// Create filters map
	filters := make(map[string]interface{})
	if status != "" {
		filters["status"] = status
	}
	if customerID != "" {
		filters["customer_id"] = customerID
	}
	page := ctx.GetInt("page")
	limit := ctx.GetInt("limit")

	orders, err := c.service.ListOrders(ctx.Request.Context(), filters, page, limit)
	if err != nil {
		c.logger.Error("Failed to list orders",
			zap.Error(err),
			zap.Any("filters", filters),
			zap.Int("page", page),
			zap.Int("limit", limit),
		)

		if err == errors.ErrInvalidStatus {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	if len(orders) == 0 {
		ctx.JSON(http.StatusOK, []domain.Order{})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles updating the status of an order
// @Summary Update order status
// @Description Updates the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param input body UpdateOrderStatusRequest true "Status update data"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders/{id}/status [put]
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var req requests.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := c.service.UpdateOrderStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		c.logger.Error("Failed to update order status",
			zap.Error(err),
			zap.String("order_id", id),
			zap.String("status", string(req.Status)),
		)

		if err == errors.ErrOrderNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
