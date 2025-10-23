package api

import (
	"math/rand"
	"order-management-ms/src/main/models/datastore"
	"order-management-ms/src/main/pkg/utils"
	"time"
)

func (req *CreateOrderRequest) ToDomain() *datastore.Order {
	return &datastore.Order{
		CustomerID: req.CustomerID,
		OrderID:    utils.GenerateOrderID(),
		Status:     datastore.StatusNew,
		Items:      req.ToOrderItem(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// assignRandomPrices assigns random prices to all items in the order
func assignRandomPrices() float64 {
	// Generate a random price between 1.00 and 100.00
	return 1.0 + float64(rand.Intn(10000))/100.0
}

func (req *CreateOrderRequest) ToOrderItem() []datastore.OrderItem {
	items := make([]datastore.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = datastore.OrderItem{
			Sku:      item.Sku,
			Quantity: item.Quantity,
			Price:    assignRandomPrices(),
		}
	}
	return items
}

func newItemFromDomain(item []datastore.OrderItem) []Items {

	items := make([]Items, len(item))
	for i, item := range item {
		items[i] = Items{
			Sku:      item.Sku,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}
	return items
}

func NewOrderResponse(order *datastore.Order) *OrderResponse {
	return &OrderResponse{
		CustomerID: order.CustomerID,
		OrderID:    order.OrderID,
		Status:     string(order.Status),
		Items:      newItemFromDomain(order.Items),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  time.Now(),
	}
}
