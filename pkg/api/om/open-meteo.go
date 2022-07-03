package om

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/romanzh1/weather-averager/pkg/api"
)

type Weather struct {
	Hourly struct {
		Time                []string  `json:"time"`
		Relativehumidity2M  []float32 `json:"relativehumidity_2m"`
		Windspeed10M        []float32 `json:"windspeed_10m"`
		Cloudcover          []float32 `json:"cloudcover"`
		Temperature2M       []float32 `json:"temperature_2m"`
		Weathercode         []float32 `json:"weathercode"`
		ApparentTemperature []float32 `json:"apparent_temperature"`
	} `json:"hourly"`
	Daily struct {
		Weathercode            []float32 `json:"weathercode"`
		Windspeed10MMax        []float32 `json:"windspeed_10m_max"`
		Temperature2MMax       []float32 `json:"temperature_2m_max"`
		Time                   []string  `json:"time"`
		ApparentTemperatureMax []float32 `json:"apparent_temperature_max"`
		PrecipitationSum       []float32 `json:"precipitation_sum"`
		Temperature2MMin       []float32 `json:"temperature_2m_min"`
		ApparentTemperatureMin []float32 `json:"apparent_temperature_min"`
	} `json:"daily"`
}

func GetCurrentDateAndHour(dataWeather *Weather) int {
	nowHour := time.Now().String()[11:13]
	nowDate := time.Now().String()[:10]

	for i, el := range dataWeather.Hourly.Time {
		if nowHour == el[11:13] && nowDate == el[:10] {
			return i
		}
	}
	return 0
}

func GetTomorrowDateAndHour(dataWeather *Weather) int {
	midnight := time.Now().AddDate(0, 0, 1).String()[:10]

	for i, el := range dataWeather.Hourly.Time {
		if midnight == el[:10] {
			return i
		}
	}
	return 0
}

func GetWeatherNow(locality string) (*Weather, int, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}
	addressAPI := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,relativehumidity_2m,apparent_temperature,weathercode,cloudcover,windspeed_10m",
		coordinates.Lat, coordinates.Lon)

	dataWeather, err := getWeather(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, GetCurrentDateAndHour(dataWeather), err
}

func GetWeatherByHour(locality string) (*Weather, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}
	addressAPI := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,relativehumidity_2m,apparent_temperature,weathercode,cloudcover,windspeed_10m",
		coordinates.Lat, coordinates.Lon)

	dataWeather, err := getWeather(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, err
}

func GetWeatherByDay(locality string) (*Weather, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(coordinates)
	addressAPI := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=weathercode,temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,precipitation_sum,windspeed_10m_max",
		coordinates.Lat, coordinates.Lon)
	addressAPI += "&timezone=Europe%" + "2FMoscow"
	fmt.Println(addressAPI)

	dataWeather, err := getWeather(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, err
}

func getWeather(addressAPI string) (*Weather, error) {
	body, err := http.Get(addressAPI)
	if err != nil {
		return nil, err
	}

	dataWeather := &Weather{}
	err = json.NewDecoder(body.Body).Decode(dataWeather)
	if err != nil {
		return nil, err
	}

	return dataWeather, err
}
