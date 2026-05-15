package config

import (
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Port                 string
	JWTSecret            string
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
	appEnv := appEnv()

	return Config{
		Port:                 envOrDefault("PORT", "8081"),
		JWTSecret:            requiredSecret("JWT_SECRET", "dev-secret", appEnv),
		InternalServiceToken: requiredSecret("INTERNAL_SERVICE_TOKEN", "dev-internal-secret", appEnv),
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
	dsn := url.URL{
		Scheme:  "postgres",
		User:    url.UserPassword(c.User, c.Password),
		Host:    net.JoinHostPort(c.Host, c.Port),
		Path:    "/" + c.Name,
		RawPath: "/" + url.PathEscape(c.Name),
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	query.Set("TimeZone", "UTC")
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func appEnv() string {
	return strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
}

func isDevelopment(appEnv string) bool {
	return appEnv == "" || appEnv == "development"
}

func requiredSecret(key, fallback, appEnv string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	if isDevelopment(appEnv) {
		return fallback
	}

	log.Fatalf("%s is required when APP_ENV=%s", key, appEnv)
	return ""
}
