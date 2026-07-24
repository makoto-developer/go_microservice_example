package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Payment      PaymentConfig
	Notification NotificationConfig
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

// NotificationConfig は通知サービスへの接続先(配達完了メールに使う)
type NotificationConfig struct {
	Host string
	Port string
}

// PaymentConfig は決済サービスへの接続先(代引きの集金確定通知に使う)
type PaymentConfig struct {
	Host string
	Port string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SHIPPING_SERVICE_PORT", "50057"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("SHIPPING_DB_HOST", "localhost"),
			Port:     getEnv("SHIPPING_DB_PORT", "5438"),
			User:     getEnv("SHIPPING_DB_USER", "postgres"),
			Password: getEnv("SHIPPING_DB_PASSWORD", "postgres"),
			DBName:   getEnv("SHIPPING_DB_NAME", "shipping_db"),
			SSLMode:  getEnv("SHIPPING_DB_SSLMODE", "disable"),
		},
		Payment: PaymentConfig{
			Host: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
			Port: getEnv("PAYMENT_SERVICE_PORT", "50056"),
		},
		Notification: NotificationConfig{
			Host: getEnv("NOTIFICATION_SERVICE_HOST", "localhost"),
			Port: getEnv("NOTIFICATION_SERVICE_PORT", "20107"),
		},
	}, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password,
		c.Database.DBName, c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
