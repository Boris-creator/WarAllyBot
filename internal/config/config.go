package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	BOT_TOKEN_KEY   = "BOT_TOKEN"
	DB_NAME_KEY     = "DB_NAME"
	DB_USER_KEY     = "DB_USER"
	DB_PASSWORD_KEY = "DB_PASSWORD"
	DB_PORT_KEY     = "DB_PORT"
)

func Load() error {
	return godotenv.Load("../../.env")
}

func Config(key string) string {
	return os.Getenv(key)
}
