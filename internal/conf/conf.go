package conf

import (
	"os"
)

type Config struct {
	JWTSecret   string
	User        string
	Pass        string
	DatabaseURL string
}

// Загрузка переменных окружения.
func Load() *Config {
	return &Config{
		JWTSecret:   getEnv("JWT_SECRET", "default-jwt-secret"),
		User:        getEnv("USERNAME", "admin"),
		Pass:        getEnv("PASSWORD", "password"),
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost port=5432 user=user password=pass dbname=db sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
