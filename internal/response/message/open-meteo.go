package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/romanzh1/weather-averager/pkg/api/om"
)

func GetWeatherOM(message string) string {
	percent := "%"
	reply := ""

	if strings.Contains(message, "сейчас") {
		locality := strings.Fields(message)
		dataWeather, hour, err := om.GetWeatherNow(locality[0]) // TODO make any order
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Sprintf("Сейчас %.1f°, ощущается как %.1f°. Влажность %.0f%s\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
			dataWeather.Hourly.Temperature2M[hour], dataWeather.Hourly.ApparentTemperature[hour], dataWeather.Hourly.Cloudcover[hour], percent,
			dataWeather.Hourly.Windspeed10M[hour], dataWeather.Hourly.Relativehumidity2M[hour], percent, getWeatherOMCondition(dataWeather.Hourly.Weathercode[hour])) //TODO change the sprintf to something else
	}

	if strings.Contains(message, "сегодня") {
		message := strings.Fields(message)
		dataWeather, err := om.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}

		reply = "Погода на 12 часов: "

		iHour := om.GetCurrentDateAndHour(dataWeather)

		for i := iHour; i < iHour+12; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Hourly.Time[i][11:16], dataWeather.Hourly.Temperature2M[i], dataWeather.Hourly.ApparentTemperature[i], dataWeather.Hourly.Cloudcover[i], percent,
				dataWeather.Hourly.Windspeed10M[i], dataWeather.Hourly.Relativehumidity2M[i], percent, getWeatherOMCondition(dataWeather.Hourly.Weathercode[i]))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, "завтра") {
		message := strings.Fields(message)
		dataWeather, err := om.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}
		reply = "Погода на 16 часов завтра: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

		iHour := om.GetTomorrowDateAndHour(dataWeather)

		for i := iHour; i < iHour+16; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Hourly.Time[i][11:16], dataWeather.Hourly.Temperature2M[i], dataWeather.Hourly.ApparentTemperature[i], dataWeather.Hourly.Cloudcover[i], percent,
				dataWeather.Hourly.Windspeed10M[i], dataWeather.Hourly.Relativehumidity2M[i], percent, getWeatherOMCondition(dataWeather.Hourly.Weathercode[i]))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, " на ") {
		message := strings.Fields(message)
		dataWeather, err := om.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}
		numberHours, err := strconv.Atoi(message[2])
		if err != nil {
			fmt.Println(err)
		}
		reply = fmt.Sprintf("Погода на %d часов", numberHours)

		iHour := om.GetCurrentDateAndHour(dataWeather)

		for i := iHour; i < iHour+numberHours; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %.0f%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				dataWeather.Hourly.Time[i][11:16], dataWeather.Hourly.Temperature2M[i], dataWeather.Hourly.ApparentTemperature[i], dataWeather.Hourly.Cloudcover[i], percent,
				dataWeather.Hourly.Windspeed10M[i], dataWeather.Hourly.Relativehumidity2M[i], percent, getWeatherOMCondition(dataWeather.Hourly.Weathercode[i]))
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

		dataWeather, err := om.GetWeatherByDay(partMessage[0])
		if err != nil {
			fmt.Println(err)
		}

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

			numDays, err = strconv.Atoi(partMessage[1])
			if err != nil {
				fmt.Println(err)
			}
			if numDays > 7 {
				numDays = 7
			}
		}

		fmt.Println(dataWeather)
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

			reply += fmt.Sprintf("\n\n%s\nФактическая: от %.1f° до %.1f°. По ощущениям: от %.1f° до %.1f°"+
				"\nСумма осадков: %.0f мм. Скорость ветра %.1f м/с, %s",
				day, dataWeather.Daily.Temperature2MMin[i], dataWeather.Daily.Temperature2MMax[i],
				dataWeather.Daily.ApparentTemperatureMin[i], dataWeather.Daily.ApparentTemperatureMax[i],
				dataWeather.Daily.PrecipitationSum[i], dataWeather.Daily.Windspeed10MMax[i], getWeatherOMCondition(dataWeather.Daily.Weathercode[i]))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	return invalidRequest
}

func getWeatherOMCondition(code float32) string {
	condition := ""

	switch code {
	case 0.0:
		condition = "ясно🌞"
	case 1.0:
		condition = "преимущественно ясно🌤"
	case 2.0:
		condition = "переменная облачность☁️"
	case 3.0:
		condition = "пасмурно☁"
	case 45.0:
		condition = "туман🌫"
	case 48.0:
		condition = "иней🌫"
	case 51.0:
		condition = "легкая морось🌧"
	case 53.0:
		condition = "умеренная морось🌧"
	case 55.0:
		condition = "морось плотной интенсивности🌧"
	case 56.0:
		condition = "слабый моросящий дождь🌧"
	case 57.0:
		condition = "моросящий дождь🌧"
	case 61.0:
		condition = "слабый дождь🌧"
	case 63.0:
		condition = "средний дождь🌧"
	case 65.0:
		condition = "сильный дождь🌧"
	case 66.0:
		condition = "слабый ледяной дождь🌧🥶"
	case 67.0:
		condition = "сильный ледяной дождь🌧🥶"
	case 71.0:
		condition = "слабый снегопад🌨"
	case 73.0:
		condition = "средний снегопад🌨"
	case 75.0:
		condition = "сильный снегопад🌨"
	case 77.0:
		condition = "снежные хлопья❄"
	case 80.0:
		condition = "не большой ливень🌧"
	case 81.0:
		condition = "средний ливень🌧"
	case 82.0:
		condition = "сильный ливень🌧"
	case 85.0:
		condition = "не большой снежный дождь❄🌧"
	case 86.0:
		condition = "сильный снежный дождь❄🌧"
	case 95.0:
		condition = "слабая гроза🌩"
	case 96.0:
		condition = "гроза с небольшим градом🌩🧊"
	case 99.0:
		condition = "гроза с сильным градом🌩🧊"
	}

	return condition
}
