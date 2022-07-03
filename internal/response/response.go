package response

import (
	"fmt"
	"strings"

	"github.com/romanzh1/weather-averager/internal/response/message"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendResponse(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) error {
	var err error
	// TODO optimize the number of if constructs and create methods for repeated code
	// TODO add notification functionality
	// TODO message design (clouds, amount of information)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		userMessage := update.Message.Text
		reply := "Не знаю, что сказать🧐Попробуй написать /help, чтобы узнать, что я могу"
		if userMessage == "" {
			reply = "Напиши что-нибудь☝️"
		}

		if !strings.Contains(userMessage, "om") && !strings.Contains(userMessage, "tom") {
			reply = message.GetWeatherOWM(userMessage)
		}

		if strings.Contains(userMessage, "tom") {
			reply = message.GetWeatherTOM(userMessage)
		} else if strings.Contains(userMessage, "om") {
			reply = message.GetWeatherOM(userMessage)
		}

		if strings.Contains(userMessage, "ave") {
			reply = message.GetWeatherAverage(userMessage)
		}
		// log.Printf("[%s] %s", update.Message.From.UserName, message)

		switch update.Message.Command() {
		case "start":
			reply = "Привет🖐. Я телеграм-бот, усредняющий погоду из различных популярных " +
				"сервисов погоды, на данный момент могу показать погоду из трёх.\n" +
				"Напиши /help, чтобы узнать, что я могу."
		case "stop":
			reply = "hammer time"
		case "help":
			reply = "Я могу показать погоду из сервисов \"OpenWeatherMap(OWM)\", \"Open-meteo(OM)\" и \"Tomorrow.io(TOM)\" " + "в твоём городе. " +
				"Для этого используй команды на клавиатуре или введи команду следующим образом:\n" +
				"	• Можайск сейчас\n" +
				"	• Тверь на 6 часов\n" +
				"	• Москва сегодня\n" +
				"	• Дубна завтра\n" +
				"	• Петербург 4 дня\n" +
				"	• Териберка неделя\n" +
				"По умолчанию используется источник OWM, если хочешь выбрать другой допиши в конце команды название источника, например \"tom\"\n" +
				"Для агрегации данных всех источников добавь в конце фразу \"ave\".\n" +
				"\nПрогноз погоды можно узнать на данный момент, на сегодня и завтра с ежечасным прогнозом.\n" +
				"Либо укажи на сколько часов показать прогноз, почасовой прогноз работает вплоть до 37 часов.\n" +
				"Также, можно посмотреть подневный прогноз на период вплоть до 7 дней.\n" +
				"\nПомимо твоих запросов, я могу автоматически присылать уведомления о погоде (На данный момент в разработке👨‍💻). Для этого введи команду так:\n" +
				"	• уведомление 12:00\n" +
				"	• уведомление 15:30 каждые 2 дня\n" +
				"	• уведомление 9:00 каждую неделю\n"

		}

		switch userMessage {
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
