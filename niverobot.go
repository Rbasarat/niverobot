package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("api_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if strings.Contains(strings.ToLower(update.Message.Text), "werkt de bot van siwa"){
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nee, Fix je bot homo!")
			//msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}

		if strings.Contains(strings.ToLower(update.Message.Text), "ping"){
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
			//msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}



	}
}