package message

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/romanzh1/weather-averager/pkg/api/om"
	"github.com/romanzh1/weather-averager/pkg/api/owm"
	"github.com/romanzh1/weather-averager/pkg/api/tom"
)

const invalidRequest = "Не корректно введены данные для прогнозирования"

func GetWeatherAverage(message string) string {
	reply := ""

	if strings.Contains(message, "сегодня") {
		message := strings.Fields(message)[0]

		var errW error
		var dataWeatherTOM *tom.Weather
		var dataWeatherOM *om.Weather
		var dataWeatherOWM *owm.WeatherTwoDays

		var wg sync.WaitGroup
		wg.Add(3)
		tomSource := func(message string, hour int) {
			defer wg.Done()
			dataWeatherTOM, errW = tom.GetWeatherByHour(message, hour)
			if errW != nil {
				fmt.Println(errW)
			}
		}
		omSource := func(message string) {
			defer wg.Done()
			dataWeatherOM, errW = om.GetWeatherByHour(message)
			if errW != nil {
				fmt.Println(errW)
			}
		}
		owmSource := func(message string) {
			defer wg.Done()
			dataWeatherOWM, errW = owm.GetWeatherByHour(message)
			if errW != nil {
				fmt.Println(errW)
			}
		}

		go tomSource(message, 12)
		go omSource(message)
		go owmSource(message)

		wg.Wait()

		reply = "Усреднённая погода на 12 часов из всех источников: " +
			buildForecastText(dataWeatherTOM, dataWeatherOM, dataWeatherOWM, 0, 12)

		return reply
	}

	if strings.Contains(message, " на ") {
		message := strings.Fields(message)

		numberHours, err := strconv.Atoi(message[2])
		if err != nil {
			fmt.Println(err)
		}

		dataWeatherTOM, err := tom.GetWeatherByHour(message[0], 12)
		if err != nil {
			fmt.Println(err)
		}
		dataWeatherOM, err := om.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}
		dataWeatherOWM, err := owm.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}

		reply = fmt.Sprintf("Усреднённая погода на %d часов из всех источников", numberHours) +
			buildForecastText(dataWeatherTOM, dataWeatherOM, dataWeatherOWM, 0, numberHours)

		return reply
	}

	return invalidRequest
}

func buildForecastText(TOM *tom.Weather, OM *om.Weather, OWM *owm.WeatherTwoDays, start, end int) string {
	forecast := ""
	percent := "%"

	iHour := om.GetCurrentDateAndHour(OM)

	for i := start; i < end; i++ {
		j := iHour + i
		forecast += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
			TOM.Data.Timelines[0].Intervals[i].StartTime.String()[11:16],
			(TOM.Data.Timelines[0].Intervals[i].Values.Temperature+OM.Hourly.Temperature2M[j]+OWM.Hourly[i].Temp)/3,
			(TOM.Data.Timelines[0].Intervals[i].Values.TemperatureApparent+OM.Hourly.ApparentTemperature[j]+OWM.Hourly[i].FeelsLike)/3,
			(TOM.Data.Timelines[0].Intervals[i].Values.Humidity+OM.Hourly.Relativehumidity2M[j]+float32(OWM.Hourly[i].Humidity))/3, percent,
			(TOM.Data.Timelines[0].Intervals[i].Values.WindSpeed+OM.Hourly.Windspeed10M[j]+OWM.Hourly[i].WindSpeed)/3,
			(TOM.Data.Timelines[0].Intervals[i].Values.CloudCover+OM.Hourly.Cloudcover[j]+OWM.Hourly[i].Pop*100)/3, percent,
			getWeatherOMCondition(OM.Hourly.Weathercode[j]))
	}

	return forecast
}
