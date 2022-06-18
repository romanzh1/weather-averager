package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/romanzh1/weather-averager/pkg/api/owm"
)

func getWeatherOWM(message string) string {
	percent := "%"
	reply := ""

	if strings.Contains(message, "сейчас") {
		locality := strings.Fields(message)
		dataWeather, err := owm.GetWeatherNow(locality[0]) // TODO make any order
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Sprintf("Сейчас %.1f°, ощущается как %.1f°. Влажность %d%s\nСкорость ветра %.1f м/с, %s",
			dataWeather.Main.Temp, dataWeather.Main.FeelsLike, dataWeather.Main.Humidity, percent,
			dataWeather.Wind.Speed, getWeatherOWMCondition(dataWeather.Weather[0].Description)) //TODO change the sprintf to something else
	}

	if strings.Contains(message, "сегодня") {
		message := strings.Fields(message)
		dataWeather, err := owm.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}

		reply = "Погода на 12 часов: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

		for i := 0; i < 12; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f° градусов, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				time.Unix(dataWeather.Hourly[i].Dt+10800, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
				percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Pop*100, percent, getWeatherOWMCondition(dataWeather.Hourly[i].Weather[0].Description))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, "завтра") {
		message := strings.Fields(message)
		dataWeather, err := owm.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}
		reply = "Погода на 16 часов завтра: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

		tomorrow := time.Now().AddDate(0, 0, 1)
		var dateLayout = "02-01-2006"
		midnight, _ := time.Parse(dateLayout, tomorrow.Format(dateLayout))
		midnightUnix := midnight.Unix() // TODO make time zone

		var indTomor int
		for i := 0; i < 48; i++ {
			if dataWeather.Hourly[i].Dt == midnightUnix {
				indTomor = i
			}
		}

		for i := indTomor; i < indTomor+16; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				time.Unix(dataWeather.Hourly[i].Dt, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
				percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Pop*100, percent, getWeatherOWMCondition(dataWeather.Hourly[i].Weather[0].Description))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	if strings.Contains(message, " на ") {
		message := strings.Fields(message)
		dataWeather, err := owm.GetWeatherByHour(message[0])
		if err != nil {
			fmt.Println(err)
		}
		numberHours, err := strconv.Atoi(message[2])
		if err != nil {
			fmt.Println(err)
		}
		reply = fmt.Sprintf("Погода на %d часов", numberHours)

		for i := 0; i < numberHours; i++ {
			reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				time.Unix(dataWeather.Hourly[i].Dt+10800, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
				percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Pop*100, percent, getWeatherOWMCondition(dataWeather.Hourly[i].Weather[0].Description))
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
		dataWeather, err := owm.GetWeatherByHour(partMessage[0])
		if err != nil {
			fmt.Println(err)
		}

		var numDays int
		if strings.Contains(message, "неделя") || strings.Contains(message, "неделю") {
			reply = "Погода на неделю: "
			numDays = len(dataWeather.Daily)
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

		days := [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
		for i := 0; i < numDays; i++ {
			day := days[time.Unix(dataWeather.Daily[i].Dt, 0).Weekday()]

			reply += fmt.Sprintf("\n\n%s\nФактическая: утром %.1f°, днём %.1f°, вечером %.1f°, ночью %.1f°"+
				"\nПо ощущениям: утром %.1f°, днём %.1f°, вечером %.1f°, ночью %.1f°."+
				"\nВлажность %d%s.\nСкорость ветра %.1f м/с. Вероятность осадков %.0f%s, %s",
				day, dataWeather.Daily[i].Temp.Morn, dataWeather.Daily[i].Temp.Day, dataWeather.Daily[i].Temp.Eve, dataWeather.Daily[i].Temp.Night,
				dataWeather.Daily[i].FeelsLike.Morn, dataWeather.Daily[i].FeelsLike.Day, dataWeather.Daily[i].FeelsLike.Eve, dataWeather.Daily[i].FeelsLike.Night, dataWeather.Daily[i].Humidity,
				percent, dataWeather.Daily[i].WindSpeed, dataWeather.Daily[i].Pop*100, percent, getWeatherOWMCondition(dataWeather.Daily[i].Weather[0].Description))
		}
		if err != nil {
			fmt.Println(err)
		}

		return reply
	}

	return "Не корректно введены данные для прогнозирования"
}

func getWeatherOWMCondition(condition string) string {
	conditionWithEmoji := ""

	switch condition {
	case "clear sky":
		conditionWithEmoji = "ясно🌞"
	case "few clouds":
		conditionWithEmoji = "преимущественно ясно🌤"
	case "scattered clouds":
		conditionWithEmoji = "переменная облачность☁️"
	case "broken clouds":
		conditionWithEmoji = "облачно с прояснениями⛅"
	case "overcast clouds":
		conditionWithEmoji = "пасмурно☁"
	case "shower rain":
		conditionWithEmoji = "ливень🌧"
	case "light rain":
		conditionWithEmoji = "небольшой дождь🌧"
	case "moderate rain":
		conditionWithEmoji = "средний дождь🌧"
	case "rain":
		conditionWithEmoji = "дождь🌧"
	case "thunderstorm":
		conditionWithEmoji = "гроза⛈"
	case "snow":
		conditionWithEmoji = "снег🌨"
	case "mist":
		conditionWithEmoji = "туман🌫"
	default:
		conditionWithEmoji = "⛅" // TODO add condition cases
	}

	return conditionWithEmoji
}
