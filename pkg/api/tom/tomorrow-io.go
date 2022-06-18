package tom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/romanzh1/weather-averager/pkg/api"
)

type Weather struct {
	Data struct {
		Timelines []struct {
			Intervals []struct {
				StartTime time.Time `json:"startTime"`
				Values    struct {
					CloudCover          float64     `json:"cloudCover"`
					Humidity            float64     `json:"humidity"`
					Temperature         float64 `json:"temperature"`
					TemperatureApparent float64 `json:"temperatureApparent"`
					WeatherCode         int     `json:"weatherCode"`
					WindSpeed           float64     `json:"windSpeed"`
				} `json:"values"`
			} `json:"intervals"`
		} `json:"timelines"`
	} `json:"data"`
}

func GetTomorrowDateAndHour(dataWeather *Weather) int {
	tomorrow := time.Now().AddDate(0, 0, 1).String()
	midnightStr := strings.Split(tomorrow, " ")[0] + "T00:00:00Z"

	for i, el := range dataWeather.Data.Timelines[0].Intervals{
		if midnightStr == convertTimeToTimestampz(el.StartTime){
			return i
		}
	}
	return 0
}

func convertTimeToTimestampz(time time.Time) string {
	strTime := strings.Split(time.String(), " ")
	return strTime[0] + "T" + strings.Split(strTime[1], ".")[0] + "Z"
}

func GetWeatherNow(locality string) (*Weather, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	addressAPI := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%f,%f&fields=temperature,temperatureApparent,cloudCover,windSpeed,humidity,weatherCode&startTime=%s&endTime=%s&timesteps=current&units=metric&timezone=UTC&apikey=%s",
		coordinates.Lat, coordinates.Lon, convertTimeToTimestampz(startTime), convertTimeToTimestampz(endTime), os.Getenv("TOMORROW_TOKEN"))
	
	dataWeather, err := getWeather(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, err
}

func GetWeatherByHour(locality string, countHour int) (*Weather, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Hour * time.Duration(countHour)) // move the time calculation here
	addressAPI := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%f,%f&fields=temperature,temperatureApparent,cloudCover,windSpeed,humidity,weatherCode&startTime=%s&endTime=%s&timesteps=1h&units=metric&timezone=UTC&apikey=%s",
		coordinates.Lat, coordinates.Lon, convertTimeToTimestampz(startTime), convertTimeToTimestampz(endTime), os.Getenv("TOMORROW_TOKEN"))

	dataWeather, err := getWeather(addressAPI)
	if err != nil {
		fmt.Println(err)
	}

	return dataWeather, err
}

func GetWeatherByDay(locality string, countDay int) (*Weather, error) {
	coordinates, err := api.GetCoordLocality(locality)
	if err != nil {
		fmt.Println(err)
	}

	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, countDay)
	addressAPI := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%f,%f&fields=temperature,temperatureApparent,cloudCover,windSpeed,humidity,weatherCode&startTime=%s&endTime=%s&timesteps=1d&units=metric&timezone=UTC&apikey=%s",
		coordinates.Lat, coordinates.Lon, convertTimeToTimestampz(startTime), convertTimeToTimestampz(endTime), os.Getenv("TOMORROW_TOKEN"))

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
