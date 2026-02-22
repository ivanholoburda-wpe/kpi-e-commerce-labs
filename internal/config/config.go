package config

import "os"

type Config struct {
	DatabaseURL   string
	TelegramToken string
	JWTSecret     []byte
	Addr          string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		Addr:          os.Getenv("ADDR"),
	}
}
