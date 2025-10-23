package mongodb

import (
	"context"
	"time"

	domain "order-management-ms/src/main/domain"
	errors "order-management-ms/src/main/pkg/customerrors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// OrderRepositoryMongoDB implements OrderRepository for MongoDB
type OrderRepositoryMongoDB struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewOrderRepository creates a new MongoDB order repository
func NewOrderRepository(db *mongo.Database, collectionName string, logger *zap.Logger) *OrderRepositoryMongoDB {
	return &OrderRepositoryMongoDB{
		collection: db.Collection(collectionName),
		logger:     logger,
	}
}

// Create saves a new order to MongoDB
func (r *OrderRepositoryMongoDB) Create(ctx context.Context, order *domain.Order) error {
	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Error("Failed to create order", zap.Error(err))
		return err
	}

	// Update the order with the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		order.ID = oid
	}

	return nil
}

// FindByID finds an order by its ID in MongoDB
func (r *OrderRepositoryMongoDB) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid order ID format", zap.Error(err), zap.String("order_id", id))
		return nil, errors.ErrOrderNotFound
	}

	var order domain.Order
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrOrderNotFound
		}
		r.logger.Error("Failed to find order", zap.Error(err), zap.String("order_id", id))
		return nil, err
	}

	return &order, nil
}

// UpdateStatus updates the status of an order in MongoDB
func (r *OrderRepositoryMongoDB) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid order ID format", zap.Error(err), zap.String("order_id", id))
		return errors.ErrOrderNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
	)

	if err != nil {
		r.logger.Error("Failed to update order status",
			zap.Error(err),
			zap.String("order_id", id),
			zap.String("status", string(status)),
		)
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrOrderNotFound
	}

	return nil
}

// FindAll finds all orders, filtered by status and customer ID
func (r *OrderRepositoryMongoDB) FindAll(ctx context.Context, filters map[string]interface{}) ([]domain.Order, error) {
	filter := bson.M{}

	// Add status filter if provided
	if status, ok := filters["status"].(string); ok {
		filter["status"] = status
	}

	// Add customer ID filter if provided
	if customerID, ok := filters["customer_id"].(string); ok {
		filter["customer_id"] = customerID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
