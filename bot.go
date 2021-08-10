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
		if update.Message.Text == ""{
			reply = "Используй только текст☝️"
		} 

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Привет🖐. Я телеграм-бот, усредняющий погоду из различных популярных " +
			"сервисов погоды. Напиши /help, чтобы узнать, что я могу"
		case "stop":
			reply = "hammer time"
		case "help":
			reply = "Я могу показать среднее арифметическое прогнозов погоды из сервисов \"Яндекс\" " +
			"в твоём городе, для этого используй команды на клавиатуре или введи команду следующим образом:\n" +
			"	погода Москва сегодня\n" +
			"	погода Тольятти неделя\n" +
			"Погоду можно спрогнозировать максимум на месяц вперёд, но учти, что это не точно" +
			"Также, есть возможность показать погоду из определённого сервиса, для этого используй уже " + 
			"перечисленные названия и введи команду так:\n" +
			"	Яндекс погода Москва сегодня"
		}

		switch update.Message.Text{
		case "Привет":
			reply = "Привет🖐. Я телеграм-бот, усредняющий погоду из различных популярных " +
			"сервисов погоды. Напиши /help, чтобы узнать, что я могу"
		case "hello":
			reply = "world😜"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
