package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/romanzh1/weather-averager/pkg/api/tom"
)

func GetWeatherTOM(message string) string {
	percent := "%"
	reply := ""

	if strings.Contains(message, "сейчас") {
		locality := strings.Fields(message)
		dataWeather, err := tom.GetWeatherNow(locality[0]) // TODO make any order
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Sprintf("Сейчас %.1f°, ощущается как %.1f°. Влажность %.0f%s\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
			dataWeather.Data.Timelines[0].Intervals[0].Values.Temperature, dataWeather.Data.Timelines[0].Intervals[0].Values.TemperatureApparent,
			dataWeather.Data.Timelines[0].Intervals[0].Values.CloudCover, percent, dataWeather.Data.Timelines[0].Intervals[0].Values.WindSpeed,
			dataWeather.Data.Timelines[0].Intervals[0].Values.Humidity, percent, getWeatherTOMCondition(dataWeather.Data.Timelines[0].Intervals[0].Values.WeatherCode)) //TODO change the sprintf to something else
	}

	if strings.Contains(message, "сегодня") {
		message := strings.Fields(message)
		dataWeather, err := tom.GetWeatherByHour(message[0], 12)
		if err != nil {
			fmt.Println(err)
		}

		reply = "Погода на 12 часов: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

		for i := 0; i < 12; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Data.Timelines[0].Intervals[i].StartTime.String()[11:16], dataWeather.Data.Timelines[0].Intervals[i].Values.Temperature, dataWeather.Data.Timelines[0].Intervals[i].Values.TemperatureApparent,
				dataWeather.Data.Timelines[0].Intervals[i].Values.CloudCover, percent, dataWeather.Data.Timelines[0].Intervals[i].Values.WindSpeed,
				dataWeather.Data.Timelines[0].Intervals[i].Values.Humidity, percent, getWeatherTOMCondition(dataWeather.Data.Timelines[0].Intervals[i].Values.WeatherCode))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, "завтра") {
		message := strings.Fields(message)
		dataWeather, err := tom.GetWeatherByHour(message[0], 41) // TODO calculate the number of hours
		if err != nil {
			fmt.Println(err)
		}
		reply = "Погода на 16 часов завтра: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

		iHour := tom.GetTomorrowDateAndHour(dataWeather)

		for i := iHour; i < iHour+16; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Data.Timelines[0].Intervals[i].StartTime.String()[11:16], dataWeather.Data.Timelines[0].Intervals[i].Values.Temperature, dataWeather.Data.Timelines[0].Intervals[i].Values.TemperatureApparent,
				dataWeather.Data.Timelines[0].Intervals[i].Values.CloudCover, percent, dataWeather.Data.Timelines[0].Intervals[i].Values.WindSpeed,
				dataWeather.Data.Timelines[0].Intervals[i].Values.Humidity, percent, getWeatherTOMCondition(dataWeather.Data.Timelines[0].Intervals[i].Values.WeatherCode))
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
		dataWeather, err := tom.GetWeatherByHour(message[0], numberHours)
		if err != nil {
			fmt.Println(err)
		}

		reply = fmt.Sprintf("Погода на %d часов", numberHours)

		for i := 0; i < numberHours; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Data.Timelines[0].Intervals[i].StartTime.String()[11:16], dataWeather.Data.Timelines[0].Intervals[i].Values.Temperature, dataWeather.Data.Timelines[0].Intervals[i].Values.TemperatureApparent,
				dataWeather.Data.Timelines[0].Intervals[i].Values.CloudCover, percent, dataWeather.Data.Timelines[0].Intervals[i].Values.WindSpeed,
				dataWeather.Data.Timelines[0].Intervals[i].Values.Humidity, percent, getWeatherTOMCondition(dataWeather.Data.Timelines[0].Intervals[i].Values.WeatherCode))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, " дня") || strings.Contains(message, " дней") ||
		strings.Contains(message, " день") ||
		strings.Contains(message, "неделя") || strings.Contains(message, "неделю") {
		partMessage := strings.Fields(message)

		var numDays int
		if strings.Contains(message, "неделя") || strings.Contains(message, "неделю") {
			reply = "Погода на неделю: "
			numDays = 7
		} else {
			if strings.Contains(message, " дня") {
				reply = "Погода на " + partMessage[1] + " дня:"
			}
			if strings.Contains(message, " дней") {
				reply = "Погода на " + partMessage[1] + " дней:"
			}
			if strings.Contains(message, " день") {
				reply = "Погода на " + partMessage[1] + " день:"
			}

			num, err := strconv.Atoi(partMessage[1])
			if err != nil {
				fmt.Println(err)
			}
			numDays = num

			if numDays > 7 {
				numDays = 7
			}

		}

		dataWeather, err := tom.GetWeatherByDay(partMessage[0], numDays)
		if err != nil {
			fmt.Println(err)
		}

		days := map[string]string{
			"Monday":    "Понедельник",
			"Tuesday":   "Вторник",
			"Wednesday": "Среда",
			"Thursday":  "Четверг",
			"Friday":    "Пятница",
			"Saturday":  "Суббота",
			"Sunday":    "Воскресенье",
		}
		for i := 0; i < numDays; i++ {
			day := days[time.Now().AddDate(0, 0, i).Weekday().String()]

			reply += fmt.Sprintf("\n\n%s\nФактическая температура %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				day, dataWeather.Data.Timelines[0].Intervals[i].Values.Temperature, dataWeather.Data.Timelines[0].Intervals[i].Values.TemperatureApparent,
				dataWeather.Data.Timelines[0].Intervals[i].Values.CloudCover, percent, dataWeather.Data.Timelines[0].Intervals[i].Values.WindSpeed,
				dataWeather.Data.Timelines[0].Intervals[i].Values.Humidity, percent, getWeatherTOMCondition(dataWeather.Data.Timelines[0].Intervals[i].Values.WeatherCode))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	return invalidRequest
}

func getWeatherTOMCondition(code int) string { // TODO add condition cases
	condition := ""
	flag := true

	switch flag {
	case code == 1000:
		condition = "ясно🌞"
	case code == 1100:
		condition = "преимущественно ясно🌤"
	case code == 1101 || code == 1102:
		condition = "переменная облачность☁️"
	case code == 1001:
		condition = "облачно☁"
	case code == 2000 || code == 2100:
		condition = "туман🌫"
	case code == 4000:
		condition = "иней🌫"
	case code == 4200:
		condition = "слабый дождь🌧"
	case code == 4001:
		condition = "дождь🌧"
	case code == 4201:
		condition = "сильный дождь🌧"
	case code == 6200:
		condition = "ледяной дождь🌧🥶"
	case code == 6201:
		condition = "сильный ледяной дождь🌧🥶"
	case code == 5100:
		condition = "слабый снегопад🌨"
	case code == 5000:
		condition = "снегопад🌨"
	case code == 5101:
		condition = "сильный снегопад🌨"
	case code == 8000:
		condition = "гроза🌩"
	case code == 7102 || code == 7000:
		condition = "гроза с градом🌩🧊"
	case code == 7101:
		condition = "гроза с сильным градом🌩🧊"
	default:
		condition = "🌤"
	}

	return condition
}
