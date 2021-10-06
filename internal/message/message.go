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

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		reply := "Не знаю, что сказать🧐Попробуй написать /help, чтобы узнать, что я могу"
		if update.Message.Text == "" {
			reply = "Используй только текст☝️"
		}

		if strings.Contains(update.Message.Text, "сейчас") {
			locality := strings.Fields(update.Message.Text)
			dataWeather, err := api.GetWeatherNow(locality[0]) // TODO make any order
			if err != nil {
				fmt.Println(err)
			}
			reply = fmt.Sprintf("Сейчас %.1f° градусов, ощущается как %.1f ⛅. Влажность %d%s\nСкорость ветра %.1f м/с, %s",
				dataWeather.Main.Temp, dataWeather.Main.FeelsLike, dataWeather.Main.Humidity, percent,
				dataWeather.Wind.Speed, dataWeather.Weather[0].Description) //TODO change the sprintf to something else
		}

		if strings.Contains(update.Message.Text, "сегодня") {
			message := strings.Fields(update.Message.Text)
			dataWeather, err := api.GetWeatherByHour(message[0]) // TODO send weather emoji
			if err != nil {
				fmt.Println(err)
			}
			reply = "Погода на 12 часов: " //TODO add a response with a proposal to send a weather forecast for another 12 hours

			for i := 0; i < 12; i++ {
				reply += fmt.Sprintf("\n\n%s %.1f° градусов, ощущается как %.1f. Влажность %d%s.\nСкорость ветра %.1f м/с, %s⛅",
					time.Unix(dataWeather.Hourly[i].Dt, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
					percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Weather[0].Description)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		if strings.Contains(update.Message.Text, "на") {
			message := strings.Fields(update.Message.Text)
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
				reply += fmt.Sprintf("\n\n%s %.1f° градусов, ощущается как %.1f. Влажность %d%s.\nСкорость ветра %.1f м/с, %s⛅",
					time.Unix(dataWeather.Hourly[i].Dt, 0).Format("15:04"), dataWeather.Hourly[i].Temp, dataWeather.Hourly[i].FeelsLike, dataWeather.Hourly[i].Humidity,
					percent, dataWeather.Hourly[i].WindSpeed, dataWeather.Hourly[i].Weather[0].Description)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

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
				"\nПомимо твоих запросов, я могу автоматически присылать уведомления о погоде. Для этого введи команду так:\n" +
				"	• уведомление 12:00\n" +
				"	• уведомление 15:30 каждые 2 дня\n" +
				"	• уведомление 9:00 каждую неделю\n"

		}

		switch update.Message.Text {
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
