package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherApi struct {
	token      string
	apiAddress string
}

type Weather struct {
	Description string `json:"description"`
}

type MainWeather struct {
	Temp     float64 `json:"temp"`
	FeelLike float64 `json:"feels_like"`
}

type WeatherData struct {
	MainWeather MainWeather `json:"main"`
	Weather     []Weather   `json:"weather"`
}

func NewWeatherApi(token string, apiAddress string) *WeatherApi {
	return &WeatherApi{
		token:      token,
		apiAddress: apiAddress,
	}
}

func MakeRequest(longitude float32, latitude float32) (map[string]string, error) {
	weatherData := WeatherData{}

	weatherMap := make(map[string]string)

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=b8561745cac51101715b64260d9d06d5")

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, err
	}

	weatherMap["real_temperature"] = fmt.Sprintf("%v", weatherData.MainWeather.Temp)
	weatherMap["feels_like"] = fmt.Sprintf("%v", weatherData.MainWeather.FeelLike)

	if len(weatherData.Weather) != 0 {
		weatherMap["description"] = fmt.Sprintf("%v", weatherData.Weather[0].Description)
	}

	return weatherMap, nil
}
