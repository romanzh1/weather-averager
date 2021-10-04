package message

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/romanzh1/weather-averager/pkg/api"
)

func SendResponse(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) error {
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
			fmt.Println(locality)
			weather, err := api.GetWeatherNow(locality[0])
			reply = "Погода на сегодня " + weather
			if err != nil {
				return err
			}
		}

		if strings.Contains(update.Message.Text, "завтра") {

		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

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
				"	• Москва сегодня\n" +
				"	• Дубна завтра\n" +
				"	• Петербург 4 дня\n" +
				"	• Териберка неделя\n" +
				"Прогноз погоды можно узнать на данный момент, на сегодня и завтра с ежечасным прогнозом или на период вплоть до 7 дней " +
				"с ежедневным прогнозом. Прогноз на сегодня также покажет состояние погоды на текущий момент.\n" +
				"Также, я могу сделать автоматически присылать уведомления о погоде. Для этого введи команду так:\n" +
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
	return nil
}
