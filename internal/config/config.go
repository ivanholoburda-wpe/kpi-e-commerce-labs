package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL   string
	TelegramToken string
	JWTSecret     []byte
	Addr          string
}

func NewConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			envOrDefault("DB_USER", "postgres"),
			envOrDefault("DB_PASSWORD", "postgres"),
			envOrDefault("DB_HOST", "localhost"),
			envOrDefault("DB_PORT", "5432"),
			envOrDefault("DB_NAME", "lab2"),
			envOrDefault("DB_SSLMODE", "disable"),
		)
	}

	return &Config{
		DatabaseURL:   dbURL,
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		Addr:          envOrDefault("ADDR", ":8080"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
