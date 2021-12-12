package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const webHook string = "https://ancient-forest-54212.herokuapp.com/"

func main() {
	port := os.Getenv("PORT")
	botToken := os.Getenv("TOKEN")

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("An error has occured while creating a bot: ", err)
	}
	log.Println("Bot created successfully")
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webHook))
	if err != nil {
		log.Fatalf("An error has occured while setting up a webhook %v: error: %v", webHook, err)
	}
	log.Println("Webhook was set successfully")

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/")
	for update := range updates {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)); err != nil {
			log.Printf("%+v\n", update)
		}
	}
}
