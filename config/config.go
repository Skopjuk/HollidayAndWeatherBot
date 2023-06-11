package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Token                string
	LogLevel             string
	BotDebug             bool
	HolidayApiToken      string
	HolidayApiUrlAddress string
	WeatherApiToken      string
	WeatherApiUrlAddress string
	MongoUrl             string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	token := os.Getenv("TOKEN")
	logLevel := os.Getenv("LOG_LEVEL")
	botDebug := os.Getenv("BOT_DEBUG")
	holidayApiToken := os.Getenv("HOLIDAY_KEY")
	holidayApiUrl := os.Getenv("HOLIDAY_API_ADDRESS")
	weatherApiToken := os.Getenv("HOLIDAY_API_TOKEN")
	weatherApiUrl := os.Getenv("WEATHER_API_ADDRESS")
	mongoUrl := os.Getenv("MONGO_URL")

	botDebugBool, err := strconv.ParseBool(botDebug)
	if err != nil {
		return nil, err
	}
	return &Config{
		Token:                token,
		LogLevel:             logLevel,
		BotDebug:             botDebugBool,
		HolidayApiToken:      holidayApiToken,
		HolidayApiUrlAddress: holidayApiUrl,
		WeatherApiToken:      weatherApiToken,
		WeatherApiUrlAddress: weatherApiUrl,
		MongoUrl:             mongoUrl,
	}, nil
}
