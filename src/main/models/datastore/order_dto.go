package datastore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusInProgress OrderStatus = "IN_PROGRESS"
	StatusDelivered  OrderStatus = "DELIVERED"
	StatusCancelled  OrderStatus = "CANCELLED"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID    string             `bson:"order_id" json:"order_id"`
	CustomerID string             `bson:"customer_id" json:"customer_id"`
	Status     OrderStatus        `bson:"status" json:"status"`
	Items      []OrderItem        `bson:"items" json:"items"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	Sku      string  `bson:"sku" json:"sku"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
}

func IsValidStatus(status OrderStatus) bool {
	switch status {
	case StatusNew, StatusInProgress, StatusDelivered, StatusCancelled:
		return true
	default:
		return false
	}
}
