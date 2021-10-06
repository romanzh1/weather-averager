package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Weather []struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherNow struct {
	Main struct {
		Dt        int64   `json:"dt"`
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather

	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Cod int `json:"cod"`
}

type WeatherTwoDays struct {
	Hourly []struct {
		Dt         int64   `json:"dt"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		Weather
	} `json:"hourly"`
}

func GetWeatherNow(locality string) (*WeatherNow, error) {
	addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=ru", locality, os.Getenv("OWM_TOKEN"))
	body, err := getBody(addressAPI)
	if err != nil {
		fmt.Println(err)
	}
	defer body.Close()

	dataWeather := &WeatherNow{}
	err = json.NewDecoder(body).Decode(dataWeather)
	if err != nil {
		fmt.Println(err)
	}
	return dataWeather, err
}

func GetWeatherByHour(locality string) (*WeatherTwoDays, error) {
	coordinates, err := GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}
	addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&units=metric&appid=%s&lang=ru",
		coordinates.Lat, coordinates.Lon, os.Getenv("OWM_TOKEN"))
	body, err := getBody(addressAPI)
	if err != nil {
		fmt.Println(err)
	}
	defer body.Close()

	dataWeather := &WeatherTwoDays{}
	err = json.NewDecoder(body).Decode(dataWeather)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, err
}

func getBody(addressAPI string) (io.ReadCloser, error) {
	resp, err := http.Get(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return resp.Body, err
}
