package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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
		if strings.Contains(err.Error(), "json: cannot unmarshal string into Go struct field") {
			coordinates, err := GetCoordLocality(locality)
			if err != nil {
				fmt.Println(err)
			}
			addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s&lang=ru",
				coordinates.Lat, coordinates.Lon, os.Getenv("OWM_TOKEN"))
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
