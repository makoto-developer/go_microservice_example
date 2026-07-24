package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Payment  PaymentServiceConfig
}

// PaymentServiceConfig は決済サービス(payment)への接続先。
type PaymentServiceConfig struct {
	Host string
	Port string
}

func (c *PaymentServiceConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("ORDER_SERVICE_PORT", "50055"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("ORDER_DB_HOST", "localhost"),
			Port:     getEnv("ORDER_DB_PORT", "5436"),
			User:     getEnv("ORDER_DB_USER", "order_user"),
			Password: getEnv("ORDER_DB_PASSWORD", "order_password"),
			DBName:   getEnv("ORDER_DB_NAME", "order_db"),
			SSLMode:  getEnv("ORDER_DB_SSLMODE", "disable"),
		},
		Payment: PaymentServiceConfig{
			Host: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
			Port: getEnv("PAYMENT_SERVICE_PORT", "50056"),
		},
	}, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
