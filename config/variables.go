package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseName      string
	Passwd            string
	User              string
	Host              string
	Port              string
	ExecuteInterval   uint64
	MaxConsumers      int
	ProducerHost      string
	ProducerQueueName string
}

func LoadEnv() Config {
	executeInterval, _ := strconv.ParseUint(os.Getenv("EXECUTE_INTERVAL"), 10, 64)
	maxConsumers, _ := strconv.ParseUint(os.Getenv("MAX_CONSUMERS"), 10, 64)
	return Config{
		DatabaseName:      os.Getenv("DATABASE_NAME"),
		Passwd:            os.Getenv("DATABASE_PASSWORD"),
		User:              os.Getenv("DATABASE_USER"),
		Host:              os.Getenv("DATABASE_HOST"),
		Port:              os.Getenv("DATABASE_PORT"),
		ExecuteInterval:   executeInterval,
		MaxConsumers:      int(maxConsumers),
		ProducerHost:      os.Getenv("PRODUCER_HOST"),
		ProducerQueueName: os.Getenv("PRODUCER_QUEUE_NAME"),
	}
}

func SetEnvExample() {
	os.Setenv("DATABASE_NAME", "example-db" /* , "database_name" */)
	os.Setenv("DATABASE_PASSWORD", "example-password" /* , "oi" */)
	os.Setenv("DATABASE_USER", "root")
	os.Setenv("DATABASE_HOST", "127.0.0.1")
	os.Setenv("DATABASE_PORT", "3306")
	os.Setenv("EXECUTE_INTERVAL", "10")
	os.Setenv("MAX_CONSUMERS", "3")
	os.Setenv("PRODUCER_HOST", "amqp://admin:admin@localhost:5672/")
	os.Setenv("PRODUCER_QUEUE_NAME", "dumont")
}
