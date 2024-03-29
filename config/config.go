package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	TickerMinutes        int
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Warning(err)
	}

	token := os.Getenv("TOKEN")
	logLevel := os.Getenv("LOG_LEVEL")
	botDebug := os.Getenv("BOT_DEBUG")
	holidayApiToken := os.Getenv("HOLIDAY_KEY")
	holidayApiUrl := os.Getenv("HOLIDAY_API_ADDRESS")
	weatherApiToken := os.Getenv("HOLIDAY_API_TOKEN")
	weatherApiUrl := os.Getenv("WEATHER_API_ADDRESS")
	mongoUrl := os.Getenv("MONGO_URL")
	tickerMinutes := os.Getenv("TICKER_TIME")
	tickerMinutesInInt, err := strconv.Atoi(tickerMinutes)

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
		TickerMinutes:        tickerMinutesInInt,
	}, nil
}
