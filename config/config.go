package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Token      string
	LogLevel   string
	BotDebug   bool
	HolidayAPI string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	token := os.Getenv("TOKEN")
	logLevel := os.Getenv("LOG_LEVEL")
	botDebug := os.Getenv("BOT_DEBUG")
	holidayAPI := os.Getenv("HOLIDAY_KEY")
	botDebugBool, err := strconv.ParseBool(botDebug)
	if err != nil {
		return nil, err
	}
	return &Config{
		Token:      token,
		LogLevel:   logLevel,
		BotDebug:   botDebugBool,
		HolidayAPI: holidayAPI,
	}, nil
}
