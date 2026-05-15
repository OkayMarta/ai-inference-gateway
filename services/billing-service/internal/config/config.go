package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                 string
	InternalServiceToken string
	FrontendOrigin       string
	DB                   DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() Config {
	return Config{
		Port:                 envOrDefault("PORT", "8081"),
		InternalServiceToken: envOrDefault("INTERNAL_SERVICE_TOKEN", "dev-internal-secret"),
		FrontendOrigin:       envOrDefault("FRONTEND_ORIGIN", "http://localhost:5173"),
		DB: DBConfig{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     envOrDefault("DB_PORT", "5432"),
			User:     envOrDefault("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     envOrDefault("DB_NAME", "billing_db"),
			SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
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
