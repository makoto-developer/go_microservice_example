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
			Port: getEnv("SEARCH_SERVICE_PORT", "20110"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("SEARCH_DB_HOST", "localhost"),
			Port:     getEnv("SEARCH_DB_PORT", "5442"),
			User:     getEnv("SEARCH_DB_USER", "postgres"),
			Password: getEnv("SEARCH_DB_PASSWORD", "postgres"),
			DBName:   getEnv("SEARCH_DB_NAME", "search_db"),
			SSLMode:  getEnv("SEARCH_DB_SSLMODE", "disable"),
		},
	}, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
