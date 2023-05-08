package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Token string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	token := os.Getenv("TOKEN")
	return &Config{Token: token}, nil
}
