package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Port                 string
	JWTSecret            string
	InternalServiceToken string
	BillingServiceURL    string
	TaskServiceURL       string
	FrontendOrigin       string
}

func Load() Config {
	appEnv := appEnv()

	return Config{
		Port:                 envOrDefault("PORT", "8080"),
		JWTSecret:            requiredSecret("JWT_SECRET", "dev-secret", appEnv),
		InternalServiceToken: requiredSecret("INTERNAL_SERVICE_TOKEN", "dev-internal-secret", appEnv),
		BillingServiceURL:    envOrDefault("BILLING_SERVICE_URL", "http://localhost:8081"),
		TaskServiceURL:       envOrDefault("TASK_SERVICE_URL", "http://localhost:8082"),
		FrontendOrigin:       envOrDefault("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
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
