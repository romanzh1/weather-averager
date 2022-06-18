package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/romanzh1/weather-averager/pkg/api/om"
	"github.com/romanzh1/weather-averager/pkg/api/owm"
	"github.com/romanzh1/weather-averager/pkg/api/tom"
)

func getWeatherAverage(message string) string {
	percent := "%"
	reply := ""

	if strings.Contains(message, "сегодня") {
		message := strings.Fields(message)

		dataWeatherTOM, err := tom.GetWeatherByHour(message[0], 12) //TODO add go routine
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
		iHour := om.GetCurrentDateAndHour(dataWeatherOM)

		reply = "Усреднённая погода на 12 часов из всех источников: "

		for i := 0; i < 12; i++ {
			j := iHour + i
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeatherTOM.Data.Timelines[0].Intervals[i].StartTime.String()[11:16], 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.Temperature + dataWeatherOM.Hourly.Temperature2M[j] + dataWeatherOWM.Hourly[i].Temp)/3, 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.TemperatureApparent + dataWeatherOM.Hourly.ApparentTemperature[j] + dataWeatherOWM.Hourly[i].FeelsLike)/3,
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.Humidity + float64(dataWeatherOM.Hourly.Relativehumidity2M[j] + dataWeatherOWM.Hourly[i].Humidity))/3, percent, 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.WindSpeed + dataWeatherOM.Hourly.Windspeed10M[j]+ dataWeatherOWM.Hourly[i].WindSpeed)/3,
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.CloudCover + float64(dataWeatherOM.Hourly.Cloudcover[j]) + float64(dataWeatherOWM.Hourly[i].Pop*100))/3, percent, 
				getWeatherOMCondition(dataWeatherOM.Hourly.Weathercode[j]))
		}
		if err != nil {
			fmt.Println(err)
		}

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
		iHour := om.GetCurrentDateAndHour(dataWeatherOM)

		reply = fmt.Sprintf("Усреднённая погода на %d часов из всех источников", numberHours)

		for i := 0; i < numberHours; i++ {
			j := iHour + i
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeatherTOM.Data.Timelines[0].Intervals[i].StartTime.String()[11:16], 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.Temperature + dataWeatherOM.Hourly.Temperature2M[j] + dataWeatherOWM.Hourly[i].Temp)/3, 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.TemperatureApparent + dataWeatherOM.Hourly.ApparentTemperature[j] + dataWeatherOWM.Hourly[i].FeelsLike)/3,
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.Humidity + float64(dataWeatherOM.Hourly.Relativehumidity2M[j] + dataWeatherOWM.Hourly[i].Humidity))/3, percent, 
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.WindSpeed + dataWeatherOM.Hourly.Windspeed10M[j]+ dataWeatherOWM.Hourly[i].WindSpeed)/3,
				(dataWeatherTOM.Data.Timelines[0].Intervals[i].Values.CloudCover + float64(dataWeatherOM.Hourly.Cloudcover[j]) + float64(dataWeatherOWM.Hourly[i].Pop*100))/3, percent, 
				getWeatherOMCondition(dataWeatherOM.Hourly.Weathercode[j]))
		}

		return reply
	}

	return "Не корректно введены данные для прогнозирования"
}