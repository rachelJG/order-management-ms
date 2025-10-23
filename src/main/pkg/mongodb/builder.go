package mongodb

import (
	"context"
	"fmt"
	"order-management-ms/src/main/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitMongoDB(cfg *config.Config, logger *zap.Logger) (*mongo.Client, error) {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Database, cfg.MongoDB.Username)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	//Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to MongoDB")
	return client, nil
}
