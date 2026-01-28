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
			Port: getEnv("CHAT_SERVICE_PORT", "20109"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("CHAT_DB_HOST", "localhost"),
			Port:     getEnv("CHAT_DB_PORT", "5441"),
			User:     getEnv("CHAT_DB_USER", "postgres"),
			Password: getEnv("CHAT_DB_PASSWORD", "postgres"),
			DBName:   getEnv("CHAT_DB_NAME", "chat_db"),
			SSLMode:  getEnv("CHAT_DB_SSLMODE", "disable"),
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
