// En src/main/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-management-ms/src/main/config"
	ordercontroller "order-management-ms/src/main/controllers"
	"order-management-ms/src/main/pkg/api"
	"order-management-ms/src/main/pkg/cache"
	"order-management-ms/src/main/pkg/kafka"
	"order-management-ms/src/main/pkg/mongodb"
	mongodbrepo "order-management-ms/src/main/repositories/mongodb"
	orderservice "order-management-ms/src/main/services/orders"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// Configure logger based on environment and log level
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatal("Error creating logger: ", err)
	}

	// Initialize MongoDB
	mongoClient, err := mongodb.InitMongoDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Initialize Redis
	redisClient := cache.InitRedis(cfg, logger)
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("Failed to close Redis connection", zap.Error(err))
		}
	}()

	// Initialize Kafka
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
	defer kafkaProducer.Close()

	// Initialize repositories
	orderRepo := mongodbrepo.NewOrderRepository(
		mongoClient.Database(cfg.MongoDB.Database),
		cfg.MongoDB.Collection,
		logger,
	)
	cacheRepo := cache.NewRedisRepository(redisClient, logger)

	// Initialize services
	orderService := orderservice.NewOrderService(orderRepo, logger, cacheRepo, kafkaProducer)

	// Initialize controllers
	orderCtrl := ordercontroller.NewOrderController(orderService, logger)

	// Run server
	runServer(cfg, orderCtrl, logger)
}

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	var logConfig zap.Config

	// Use development config in non-production environments
	if cfg.Environment == "development" {
		logConfig = zap.NewDevelopmentConfig()
	} else {
		logConfig = zap.NewProductionConfig()
	}

	// Set log level
	switch cfg.LogLevel {
	case "debug":
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Build the logger
	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	logger.Info("Logger initialized",
		zap.String("environment", cfg.Environment),
		zap.String("logLevel", cfg.LogLevel),
	)

	return logger, nil
}

func runServer(cfg *config.Config, orderCtrl *ordercontroller.OrderController, logger *zap.Logger) {
	// Configure router
	router := api.SetupRouter(orderCtrl, logger)

	// Configure HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Run the server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for signal to stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Controlled shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
