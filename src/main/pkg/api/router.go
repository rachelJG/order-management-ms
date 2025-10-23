package api

import (
	"net/http"
	"time"

	ordercontroller "order-management-ms/src/main/controllers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter configure the router
func SetupRouter(orderCtrl *ordercontroller.OrderController, logger *zap.Logger) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(logger))
	r.Use(jsonContentTypeMiddleware())

	// API v1 routes
	setupV1Routes(r, orderCtrl)

	return r
}

// setupV1Routes configure the routes for API v1
func setupV1Routes(r *gin.Engine, orderCtrl *ordercontroller.OrderController) {
	v1 := r.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", healthCheck)

		// Order routes
		ordersGroup := v1.Group("/orders")
		{
			ordersGroup.POST("", orderCtrl.CreateOrder)
			ordersGroup.GET("/:id", orderCtrl.GetOrder)
			ordersGroup.GET("", orderCtrl.ListOrders)
			ordersGroup.PATCH("/:id/status", orderCtrl.UpdateOrderStatus)
		}
	}
}

// healthCheck handle the health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// jsonContentTypeMiddleware ensure that all responses are JSON
func jsonContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// loggingMiddleware log the request information
func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request information
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process the request
		c.Next()

		// Log the response information
		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("Request processed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		)
	}
}
