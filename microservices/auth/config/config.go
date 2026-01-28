package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	JWTAccessSecret    string
	JWTRefreshSecret   string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:       getEnv("AUTH_SERVER_PORT", "50051"),
		DatabaseURL:      getEnv("AUTH_DATABASE_URL", ""),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
	}
	
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("AUTH_DATABASE_URL is required")
	}
	
	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
