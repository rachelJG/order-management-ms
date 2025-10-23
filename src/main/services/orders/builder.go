package orders

import (
	"order-management-ms/src/main/pkg/cache"
	"order-management-ms/src/main/pkg/kafka"
	"order-management-ms/src/main/repositories"

	"go.uber.org/zap"
)

type OrderService struct {
	repo   repositories.OrderRepository
	logger *zap.Logger
	cache  cache.Repository
	kafka  *kafka.Producer
}

func NewOrderService(repo repositories.OrderRepository, logger *zap.Logger, cache cache.Repository, kafka *kafka.Producer) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
		cache:  cache,
		kafka:  kafka,
	}
}
