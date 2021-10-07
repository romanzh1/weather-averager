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
		fmt.Println(err) // TODO on heroku crashes if put a fatal shutdown
	}
	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.ListenForWebhook("/" + bot.Token)
	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if os.Getenv("PORT") == "" {
		// long pooling (local) 	//TODO change environment variable for automatic switching
		updates, err = bot.GetUpdatesChan(u)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = message.SendResponse(updates, bot)
	if err != nil {
		fmt.Println(err)
	}
}
