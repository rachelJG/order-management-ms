package orders

import (
	kafkaDto "order-management-ms/src/main/models/kafka"
	"order-management-ms/src/main/pkg/cache"

	"order-management-ms/src/main/repositories"

	"go.uber.org/zap"
)

type OrderService struct {
	repo           repositories.OrderRepository
	logger         *zap.Logger
	cache          cache.Repository
	eventPublisher kafkaDto.EventPublisher
}

func NewOrderService(repo repositories.OrderRepository, logger *zap.Logger, cache cache.Repository, eventPublisher kafkaDto.EventPublisher) *OrderService {
	return &OrderService{
		repo:           repo,
		logger:         logger,
		cache:          cache,
		eventPublisher: eventPublisher,
	}
}
