package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"order-management-ms/src/main/controllers"
	"order-management-ms/src/main/models/api"
	dm "order-management-ms/src/main/models/datastore"
)

// Helper function to create a test router with the controller
func setupTestRouter(ctrl *controllers.OrderController) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/orders", ctrl.CreateOrder)
		api.GET("/orders/:id", ctrl.GetOrder)
		api.GET("/orders", ctrl.ListOrders)
		api.PATCH("/orders/:id/state", ctrl.UpdateOrderStatus)
	}
	return r
}

// Helper function to create a test order request
func createTestOrderRequest() *api.CreateOrderRequest {
	return &api.CreateOrderRequest{
		CustomerID: "customer-123",
		Items: []api.Items{
			{
				Sku:      "SKU-123",
				Quantity: 2,
			},
		},
	}
}

func TestCreateOrder(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		requestBody    *api.CreateOrderRequest
		setupMock      func(*mockOrderService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful order creation",
			requestBody: createTestOrderRequest(),
			setupMock: func(mockSvc *mockOrderService) {
				mockSvc.On("CreateOrder", mock.Anything, mock.Anything).
					Return(createTestOrderResponse(), nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"order_id":"order-123","customer_id":"customer-123","status":"NEW","items":[{"sku":"SKU-123","quantity":2,"price":10.99}],"created_at":"%s","updated_at":"%s"}`,
		},
		{
			name:           "invalid request body",
			requestBody:    &api.CreateOrderRequest{},
			setupMock:      func(mockSvc *mockOrderService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &mockOrderService{}
			tt.setupMock(mockSvc)

			// Create controller with mock service
			ctrl := controllers.NewOrderController(mockSvc, zap.NewNop())

			// Setup test router
			r := setupTestRouter(ctrl)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)

			// Make request
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				// Parse the response to get the actual timestamps
				var actualResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
				assert.NoError(t, err)

				// Get timestamps from response
				createdAt, hasCreatedAt := actualResponse["created_at"].(string)
				updatedAt, hasUpdatedAt := actualResponse["updated_at"].(string)

				// Only check timestamps if they exist in the response
				if hasCreatedAt && hasUpdatedAt {
					// Format the expected body with actual timestamps
					expectedBody := `{"order_id":"order-123","customer_id":"customer-123","status":"NEW","items":[{"sku":"SKU-123","quantity":2,"price":10.99}],"created_at":"` +
						createdAt + `","updated_at":"` + updatedAt + `"}`

					// Compare the JSON objects, ignoring the order of fields
					assert.JSONEq(t, expectedBody, w.Body.String())
				} else {
					// If timestamps are not present, just check the basic structure
					expectedBasic := `{"order_id":"order-123","customer_id":"customer-123","status":"NEW","items":[{"sku":"SKU-123","quantity":2,"price":10.99}]}`
					var expectedMap, actualMap map[string]interface{}
					json.Unmarshal([]byte(expectedBasic), &expectedMap)
					json.Unmarshal(w.Body.Bytes(), &actualMap)

					// Remove timestamps for comparison
					delete(actualMap, "created_at")
					delete(actualMap, "updated_at")

					assert.Equal(t, expectedMap, actualMap)
				}
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestGetOrder(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		orderID        string
		setupMock      func(*mockOrderService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "successful order retrieval",
			orderID: "order-123",
			setupMock: func(mockSvc *mockOrderService) {
				mockSvc.On("GetOrder", mock.Anything, "order-123").
					Return(createTestOrderResponse(), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "", // Will be handled specially in the test
		},
		{
			name:    "order not found",
			orderID: "nonexistent",
			setupMock: func(mockSvc *mockOrderService) {
				mockSvc.On("GetOrder", mock.Anything, "nonexistent").
					Return(nil, errors.New("order not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"order not found"}`,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &mockOrderService{}
			tt.setupMock(mockSvc)

			// Create controller with mock service
			ctrl := controllers.NewOrderController(mockSvc, zap.NewNop())

			// Setup test router
			r := setupTestRouter(ctrl)

			// Create request URL for getting a specific order
			url := "/api/v1/orders/" + tt.orderID

			// Make request
			req, _ := http.NewRequest(http.MethodGet, url, nil)

			// Execute request
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// For successful case (StatusOK), verify the response structure
			if tt.expectedStatus == http.StatusOK {
				// Parse the response to verify structure
				var responseMap map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseMap)
				assert.NoError(t, err)

				// Verify required fields exist and have correct values
				assert.Equal(t, "order-123", responseMap["order_id"])
				assert.Equal(t, "customer-123", responseMap["customer_id"])
				assert.Equal(t, "NEW", responseMap["status"])

				// Verify items
				items, ok := responseMap["items"].([]interface{})
				assert.True(t, ok && len(items) == 1)
				item := items[0].(map[string]interface{})
				assert.Equal(t, "SKU-123", item["sku"])
				assert.EqualValues(t, 2, item["quantity"])
				assert.EqualValues(t, 10.99, item["price"])
			} else if tt.expectedBody != "" {
				// For error cases, check the error message
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

// Helper function to create a test order domain model
func createTestOrder() *dm.Order {
	return &dm.Order{
		OrderID:    "order-123",
		CustomerID: "customer-123",
		Status:     dm.StatusNew,
		Items: []dm.OrderItem{
			{
				Sku:      "SKU-123",
				Quantity: 2,
				Price:    10.99,
			},
		},
	}
}

// Helper function to create a test order response
func createTestOrderResponse() *api.OrderResponse {
	order := createTestOrder()
	items := make([]api.Items, len(order.Items))
	for i, item := range order.Items {
		items[i] = api.Items{
			Sku:      item.Sku,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return &api.OrderResponse{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Items:      items,
	}
}
