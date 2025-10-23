package kafka

import (
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// EnsureTopic ensures that the topic exists in the Kafka cluster
// If the topic does not exist, it creates it

func EnsureTopic(broker []string, topic string, logger *zap.Logger) {

	if len(broker) == 0 {
		logger.Error("No brokers provided")
		return
	}
	// Connect to the first broker
	conn, err := kafka.Dial("tcp", broker[0])
	if err != nil {
		logger.Error("Failed to connect to Kafka", zap.Error(err))
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		logger.Error("Failed to get controller", zap.Error(err))
		return
	}

	controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
	if err != nil {
		logger.Error("Failed to connect to controller", zap.Error(err))
		return
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		logger.Error("Topic %s already exists or could not be created", zap.String("topic", topic), zap.Error(err))
	} else {
		logger.Info("Topic %s created successfully", zap.String("topic", topic))
	}
}
