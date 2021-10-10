package message

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/romanzh1/weather-averager/pkg/api"
)

func SendResponse(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) error {
	percent := "%"
	var err error
	// TODO optimize the number of if constructs and create methods for repeated code
	// TODO add notification functionality
	// TODO message design (clouds, amount of information)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		message := update.Message.Text
		reply := "Не знаю, что сказать🧐Попробуй написать /help, чтобы узнать, что я могу"
		if message == "" {
			reply = "Используй только текст☝️"
		}

		if strings.Contains(message, "сейчас") {
			locality := strings.Fields(message)
			dataWeather, err := api.GetWeatherNow(locality[0]) // TODO make any order
			if err != nil {
				fmt.Println(err)
			}
			reply = fmt.Sprintf("Сейчас %.1f°, ощущается как %.1f°⛅. Влажность %d%s\nСкорость ветра %.1f м/с, %s",
				dataWeather.Main.Temp, dataWeather.Main.FeelsLike, dataWeather.Main.Humidity, percent,
				dataWeather.Wind.Speed, dataWeather.Weather[0].Description) //TODO change the sprintf to something else
		}

		if strings.Contains(message, "сегодня") {
			message := strings.Fields(message)
			dataWeather, err := api.GetWeatherByHour(message[0]) // TODO send weather emoji
			if err != nil {
				fmt.Println(err)
			}

			reply = "Погода на 12 часов: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

			for i := 0; i < 12; i++ {
				reply += fmt.Sprintf("\n\n%s %.1f° градусов, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с, %s⛅Вероятность осадков %.0f%s",
					time.Unix(dataWeather.Hourly[i].Dt+10800, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
					percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Weather[0].Description, dataWeather.Hourly[i].Pop*100, percent)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		if strings.Contains(message, "завтра") {
			message := strings.Fields(message)
			dataWeather, err := api.GetWeatherByHour(message[0]) // TODO send weather emoji
			if err != nil {
				fmt.Println(err)
			}
			reply = "Погода на 12 часов завтра: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

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

			for i := indTomor; i < indTomor+12; i++ {
				reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с, %s⛅Вероятность осадков %.0f%s",
					time.Unix(dataWeather.Hourly[i].Dt, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
					percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Weather[0].Description, dataWeather.Hourly[i].Pop*100, percent)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		if strings.Contains(message, " на ") {
			message := strings.Fields(message)
			dataWeather, err := api.GetWeatherByHour(message[0]) // TODO send weather emoji
			if err != nil {
				fmt.Println(err)
			}
			numberHours, err := strconv.Atoi(message[2])
			if err != nil {
				fmt.Println(err)
			}
			reply = fmt.Sprintf("Погода на %d часов", numberHours)

			for i := 0; i < numberHours; i++ {
				reply += fmt.Sprintf("\n\n%s %.1f°, ощущается как %.1f°. Влажность %d%s.\nСкорость ветра %.1f м/с, %s⛅Вероятность осадков %.0f%s",
					time.Unix(dataWeather.Hourly[i].Dt+10800, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
					percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Weather[0].Description, dataWeather.Hourly[i].Pop*100, percent)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		if strings.Contains(message, " дня") || strings.Contains(message, " дней") ||
			strings.Contains(message, " день") ||
			strings.Contains(message, "неделя") || strings.Contains(message, "неделю") {
			partMessage := strings.Fields(message)
			dataWeather, err := api.GetWeatherByHour(partMessage[0]) // TODO send weather emoji
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
					"\nВлажность %d%s.\nСкорость ветра %.1f м/с, %s⛅ Вероятность осадков %.0f%s",
					day, dataWeather.Daily[i].Temp.Morn, dataWeather.Daily[i].Temp.Day, dataWeather.Daily[i].Temp.Eve, dataWeather.Daily[i].Temp.Night,
					dataWeather.Daily[i].FeelsLike.Morn, dataWeather.Daily[i].FeelsLike.Day, dataWeather.Daily[i].FeelsLike.Eve, dataWeather.Daily[i].FeelsLike.Night, dataWeather.Daily[i].Humidity,
					percent, dataWeather.Daily[i].WindSpeed, dataWeather.Daily[i].Weather[0].Description, dataWeather.Daily[i].Pop*100, percent)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, message)

		switch update.Message.Command() {
		case "start":
			reply = "Привет🖐. Я телеграм-бот, усредняющий погоду из различных популярных " +
				"сервисов погоды, на данный момент могу показать погоду только из одного.\n" +
				"Напиши /help, чтобы узнать, что я могу."
		case "stop":
			reply = "hammer time"
		case "help":
			reply = "Я могу показать погоду из сервиса \"OpenWeatherMap\" " + "в твоём городе. " +
				"Для этого используй команды на клавиатуре или введи команду следующим образом:\n" +
				"	• Можайск сейчас\n" +
				"	• Тверь на 6 часов\n" +
				"	• Москва сегодня\n" +
				"	• Дубна завтра\n" +
				"	• Петербург 4 дня\n" +
				"	• Териберка неделя\n" +
				"\nПрогноз погоды можно узнать на данный момент, на сегодня и завтра с ежечасным прогнозом.\n" +
				"Либо укажи на сколько часов показать прогноз, почасовой прогноз работает вплоть до 48 часов.\n" +
				"Также, можно посмотреть подневный прогноз на период вплоть до 7 дней.\n" +
				"\nПомимо твоих запросов, я могу автоматически присылать уведомления о погоде (На данный момент в разработке👨‍💻). Для этого введи команду так:\n" +
				"	• уведомление 12:00\n" +
				"	• уведомление 15:30 каждые 2 дня\n" +
				"	• уведомление 9:00 каждую неделю\n"

		}

		switch message {
		case "Привет":
			reply = "Привет🖐. Я телеграм-бот, усредняющий погоду из различных популярных " +
				"сервисов погоды, на данный момент могу показать погоду только из одного.\n" +
				"Напиши /help, чтобы узнать, что я могу."
		case "hello":
			reply = "world😜"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
	return err
}
