package repositories

import (
	"context"
	"time"

	domain "order-management-ms/src/main/models/datastore"
	errors "order-management-ms/src/main/pkg/customerrors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (r *OrderRepositoryMongoDB) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	// Update the order with the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		order.ID = oid
	}

	var createdOrder domain.Order
	err = r.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdOrder)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil
}

// FindByID finds an order by its ID in MongoDB
func (r *OrderRepositoryMongoDB) FindByID(ctx context.Context, orderID string) (*domain.Order, error) {

	var order domain.Order
	err := r.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrOrderNotFound
		}
		r.logger.Error("Failed to find order", zap.Error(err), zap.String("order_id", orderID))
		return nil, err
	}

	return &order, nil
}

// UpdateStatus updates the status of an order in MongoDB
func (r *OrderRepositoryMongoDB) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"order_id": orderID},
		update,
	)

	if err != nil {
		r.logger.Error("Failed to update order status",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.String("status", string(status)),
		)
		return err
	}

	if result.MatchedCount == 0 {
		return errors.ErrOrderNotFound
	}

	return nil
}

// List finds all orders, filtered by status and customer ID
func (r *OrderRepositoryMongoDB) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*domain.Order, error) {
	// Implement pagination
	skip := (page - 1) * limit
	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by creation date descending

	// Build the filter
	mongoFilter := bson.M{}
	if status, ok := filter["status"].(string); ok {
		mongoFilter["status"] = status
	}
	if customerID, ok := filter["customer_id"].(string); ok {
		mongoFilter["customer_id"] = customerID
	}

	cursor, err := r.collection.Find(ctx, mongoFilter, options)
	if err != nil {
		r.logger.Error("Failed to list orders", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		r.logger.Error("Failed to decode orders", zap.Error(err))
		return nil, err
	}

	return orders, nil
}
