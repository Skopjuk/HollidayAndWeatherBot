package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Token    string
	LogLevel string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	token := os.Getenv("TOKEN")
	logLevel := os.Getenv("LOG_LEVEL")
	return &Config{
		Token:    token,
		LogLevel: logLevel,
	}, nil
}
