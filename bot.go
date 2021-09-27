package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm WeatherAveragerBot!"))
}

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Для получения через long pooling. Включить для локального запуска
	// updates, err := bot.GetUpdatesChan(u)
	// Для получения через webhook. Включить для деплоя на heroku
	updates := bot.ListenForWebhook("/" + bot.Token)
	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		reply := "Не знаю, что сказать🧐Попробуй написать /help, чтобы узнать, что я могу"
		if update.Message.Text == "" {
			reply = "Используй только текст☝️"
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
				"	• Москва сегодня\n" +
				"	• Дубна завтра\n" +
				"	• Петербург 4 дня\n" +
				"	• Териберка неделя\n" +
				"Прогноз погоды можно узнать на сегодня и завтра с ежечасным прогнозом или на период вплоть до 7 дней " +
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

		bot.Send(msg)
	}
}
