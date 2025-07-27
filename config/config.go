package config

import (
	"os"
	"time"
)

type Config struct {
	Database DatabaseConfig
	Outbox   OutboxConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type OutboxConfig struct {
	PollingInterval time.Duration
	BatchSize       int
}

func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3307"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "transaction-outbox-practice"),
		},
		Outbox: OutboxConfig{
			PollingInterval: 5 * time.Second,
			BatchSize:       10,
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}