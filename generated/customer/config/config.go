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
			Port: getEnv("CUSTOMER_SERVICE_PORT", "50053"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("CUSTOMER_DB_HOST", "localhost"),
			Port:     getEnv("CUSTOMER_DB_PORT", "5434"),
			User:     getEnv("CUSTOMER_DB_USER", "customer_user"),
			Password: getEnv("CUSTOMER_DB_PASSWORD", "customer_password"),
			DBName:   getEnv("CUSTOMER_DB_NAME", "customer_db"),
			SSLMode:  getEnv("CUSTOMER_DB_SSLMODE", "disable"),
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
