package owm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/romanzh1/weather-averager/pkg/api"
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
		Temp      float32 `json:"temp"`
		FeelsLike float32 `json:"feels_like"`
		TempMin   float32 `json:"temp_min"`
		TempMax   float32 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather

	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
	Cod int `json:"cod"`
}

type WeatherTwoDays struct {
	Hourly []struct {
		Dt         int64   `json:"dt"`
		Temp       float32 `json:"temp"`
		FeelsLike  float32 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		Visibility int     `json:"visibility"`
		WindSpeed  float32 `json:"wind_speed"`
		Pop        float32 `json:"pop"`
		Weather
	} `json:"hourly"`

	Daily []struct {
		Dt        int64   `json:"dt"`
		Sunrise   int64   `json:"sunrise"`
		Sunset    int64   `json:"sunset"`
		Moonrise  int64   `json:"moonrise"`
		Moonset   int64   `json:"moonset"`
		MoonPhase float32 `json:"moon_phase"`
		Temp      struct {
			Day   float32 `json:"day"`
			Min   float32 `json:"min"`
			Max   float32 `json:"max"`
			Night float32 `json:"night"`
			Eve   float32 `json:"eve"`
			Morn  float32 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float32 `json:"day"`
			Night float32 `json:"night"`
			Eve   float32 `json:"eve"`
			Morn  float32 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		WindSpeed float32 `json:"wind_speed"`
		Weather
		Clouds int     `json:"clouds"`
		Pop    float32 `json:"pop"`
		Uvi    float32 `json:"uvi"`
		Rain   float32 `json:"rain,omitempty"`
	} `json:"daily"`
}

func GetWeatherNow(locality string) (*WeatherNow, error) {
	addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=en", locality, os.Getenv("OWM_TOKEN"))
	body, err := getBody(addressAPI)
	if err != nil {
		fmt.Println(err)
	}
	defer body.Close()

	dataWeather := &WeatherNow{}
	err = json.NewDecoder(body).Decode(dataWeather)
	if err != nil {
		if strings.Contains(err.Error(), "json: cannot unmarshal string into Go struct field") {
			coordinates, err := api.GetCoordLocality(locality)
			if err != nil {
				fmt.Println(err)
			}
			addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s&lang=en",
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

func GetWeatherByHour(locality string) (*WeatherTwoDays, error) { //hourly forecast for 2 days, daily forecast for 7 days
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}
	addressAPI := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&units=metric&appid=%s&lang=en",
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
