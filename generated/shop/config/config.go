package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SHOP_SERVICE_PORT", "20101"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("SHOP_DB_HOST", "localhost"),
			Port:     getEnv("SHOP_DB_PORT", "5433"),
			User:     getEnv("SHOP_DB_USER", "postgres"),
			Password: getEnv("SHOP_DB_PASSWORD", "postgres"),
			DBName:   getEnv("SHOP_DB_NAME", "shop_db"),
			SSLMode:  getEnv("SHOP_DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("SHOP_REDIS_HOST", "localhost"),
			Port:     getEnv("SHOP_REDIS_PORT", "6380"),
			Password: getEnv("SHOP_REDIS_PASSWORD", ""),
			DB:       0,
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "shop-images"),
			UseSSL:    false,
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

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
