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
	Main string `json:"main"`
}

type MainWeather struct {
	Temp     float64 `json:"temp"`
	FeelLike float64 `json:"feels_like"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
	Humidity float64 `json:"humidity"`
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

func (w *WeatherApi) MakeRequest(longitude float64, latitude float64) (map[string]string, error) {
	weatherData := WeatherData{}

	weatherMap := make(map[string]string)

	url := fmt.Sprintf("%sdata/2.5/weather?lat=%f&lon=%f&appid=%s", w.apiAddress, latitude, longitude, w.token)

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
	weatherMap["main"] = fmt.Sprintf("%v", weatherData.Weather[0].Main)
	weatherMap["temp_min"] = fmt.Sprintf("%v", weatherData.MainWeather.TempMin)
	weatherMap["temp_max"] = fmt.Sprintf("%v", weatherData.MainWeather.TempMax)
	weatherMap["humidity"] = fmt.Sprintf("%v", weatherData.MainWeather.Humidity)

	return weatherMap, nil
}
