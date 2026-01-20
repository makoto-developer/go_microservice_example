package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
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

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PAYMENT_SERVICE_PORT", "50056"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("PAYMENT_DB_HOST", "localhost"),
			Port:     getEnv("PAYMENT_DB_PORT", "5437"),
			User:     getEnv("PAYMENT_DB_USER", "postgres"),
			Password: getEnv("PAYMENT_DB_PASSWORD", "postgres"),
			DBName:   getEnv("PAYMENT_DB_NAME", "payment_db"),
			SSLMode:  getEnv("PAYMENT_DB_SSLMODE", "disable"),
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
