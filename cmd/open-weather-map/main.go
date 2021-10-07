package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/romanzh1/weather-averager/internal/message"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm WeatherAveragerBot!"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err) //TODO use zap for error handling
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//TODO add environment variable for automatic switching
	var updates tgbotapi.UpdatesChannel
	if os.Getenv("PORT") == "" {
		// long pooling (local)
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// getting through a webhook (deployment to heroku)
		updates = bot.ListenForWebhook("/" + bot.Token)
		http.HandleFunc("/", MainHandler)
		go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = message.SendResponse(updates, bot)
	if err != nil {
		fmt.Println(err)
	}
}
