package config

import "os"

type Config struct {
	Port              string
	JWTSecret         string
	BillingServiceURL string
	TaskServiceURL    string
	FrontendOrigin    string
}

func Load() Config {
	return Config{
		Port:              envOrDefault("PORT", "8080"),
		JWTSecret:         envOrDefault("JWT_SECRET", "dev-secret"),
		BillingServiceURL: envOrDefault("BILLING_SERVICE_URL", "http://localhost:8081"),
		TaskServiceURL:    envOrDefault("TASK_SERVICE_URL", "http://localhost:8082"),
		FrontendOrigin:    envOrDefault("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
