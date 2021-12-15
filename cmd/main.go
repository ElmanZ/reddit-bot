package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	reddit "github.com/ElmanZ/reddit_bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const webHook string = "https://elmans-reddit-bot.herokuapp.com/"

func main() {
	port := os.Getenv("PORT")
	botToken := os.Getenv("TOKEN")
	reddit.Get("https://www.reddit.com/r/golang.json")
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
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome! This is a Reditt bot."
			case "best":
				msg.Text = strings.Trim(fmt.Sprint(reddit.Get(reddit.Best+".json&limit=10")), "[]")
			case "rising":
				msg.Text = strings.Trim(fmt.Sprint(reddit.Get(reddit.Rising+".json&limit=10")), "[]")
			case "random":
				msg.Text = strings.Trim(fmt.Sprint(reddit.Get(reddit.Random+".json&limit=10")), "[]")
			case "subreddit":
				msg.Text = strings.Trim(fmt.Sprint(reddit.Get(reddit.SubReddit+".json&limit=10")), "[]")
			default:
				msg.Text = "Please use valid command!"
			}
		}
	}
}
