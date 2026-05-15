package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                 string
	OllamaURL            string
	BillingServiceURL    string
	InternalServiceToken string
	FrontendOrigin       string
	DB                   DBConfig
	Redis                RedisConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr            string
	Password        string
	DB              int
	CacheTTLSeconds int
}

func Load() Config {
	return Config{
		Port:                 envOrDefault("PORT", "8082"),
		OllamaURL:            envOrDefault("OLLAMA_URL", "http://localhost:11434"),
		BillingServiceURL:    envOrDefault("BILLING_SERVICE_URL", "http://localhost:8081"),
		InternalServiceToken: envOrDefault("INTERNAL_SERVICE_TOKEN", "dev-internal-secret"),
		FrontendOrigin:       envOrDefault("FRONTEND_ORIGIN", "http://localhost:5173"),
		DB: DBConfig{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     envOrDefault("DB_PORT", "5432"),
			User:     envOrDefault("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     envOrDefault("DB_NAME", "task_db"),
			SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:            envOrDefault("REDIS_ADDR", "localhost:6379"),
			Password:        os.Getenv("REDIS_PASSWORD"),
			DB:              envIntOrDefault("REDIS_DB", 0),
			CacheTTLSeconds: envIntOrDefault("CACHE_TTL_SECONDS", 60),
		},
	}
}

func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
