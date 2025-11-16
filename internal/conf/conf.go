package conf

import (
	"log"
	"os"
)

type Config struct {
	JWTSecret string
	User      string
	Pass      string
}

func Load() *Config {
	return &Config{
		JWTSecret: getEnv("JWT_SECRET"),
		User:      getEnv("USERNAME"),
		Pass:      getEnv("PASSWORD"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is not set or empty", key)
	}
	return value
}
