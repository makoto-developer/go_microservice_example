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
			Port: getEnv("SHOP_SERVICE_PORT", "50052"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("SHOP_DB_HOST", "localhost"),
			Port:     getEnv("SHOP_DB_PORT", "5433"),
			User:     getEnv("SHOP_DB_USER", "shop_user"),
			Password: getEnv("SHOP_DB_PASSWORD", "shop_password"),
			DBName:   getEnv("SHOP_DB_NAME", "shop_db"),
			SSLMode:  getEnv("SHOP_DB_SSLMODE", "disable"),
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
