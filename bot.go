package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
    resp.Write([]byte("Hi there! I'm DndSpellsBot!"))
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

	// Для получения через long pooling. Сейчас нам это не нужно.
	// updates, err := bot.GetUpdatesChan(u) 
	// Для получения через webhook
	updates := bot.ListenForWebhook("/" + bot.Token)

	http.HandleFunc("/", MainHandler)
    go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
