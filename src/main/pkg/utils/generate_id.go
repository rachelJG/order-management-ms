package utils

import "github.com/google/uuid"

func GenerateOrderID() string {
	shortID := uuid.New().String()[:8]
	return "ORD-" + shortID
}
