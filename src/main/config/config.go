package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server      Server
	MongoDB     MongoDB
	Redis       Redis
	Kafka       Kafka
	Environment string `envconfig:"ENVIRONMENT" default:"production"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
}

type Server struct {
	Port             string        `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout      time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout     time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s"`
	GracefulShutdown time.Duration `envconfig:"GRACEFUL_SHUTDOWN" default:"10s"`
}

type MongoDB struct {
	URI        string `envconfig:"MONGODB_URI" required:"true"`
	Database   string `envconfig:"MONGODB_DATABASE" default:"order_management"`
	Collection string `envconfig:"MONGODB_COLLECTION" default:"orders"`
}

type Redis struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

type Kafka struct {
	Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
	Topic   string   `envconfig:"KAFKA_TOPIC" default:"order_events"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Server.Port == "80" || c.Server.Port == "443"
}
