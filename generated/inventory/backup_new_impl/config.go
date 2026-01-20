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

func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("INVENTORY_SERVICE_PORT", "50054"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("INVENTORY_DB_HOST", "localhost"),
			Port:     getEnv("INVENTORY_DB_PORT", "5435"),
			User:     getEnv("INVENTORY_DB_USER", "inventory_user"),
			Password: getEnv("INVENTORY_DB_PASSWORD", "inventory_password"),
			DBName:   getEnv("INVENTORY_DB_NAME", "inventory_db"),
			SSLMode:  getEnv("INVENTORY_DB_SSLMODE", "disable"),
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
