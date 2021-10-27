package features

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"niverobot/model"
	"strings"
)

type AddKudo struct {
	kudos      model.Kudos
	users      model.Users
	kudoCounts model.KudoCounts
}

func NewAddKudo(kudosService model.Kudos, userService model.Users, kudoCountService model.KudoCounts) AddKudo {
	return AddKudo{kudos: kudosService, users: userService, kudoCounts: kudoCountService}
}

func (k AddKudo) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI, history model.MessageHistory) {
	var sendMsg tgbotapi.MessageConfig
	var lastMessage *tgbotapi.Message
	for i := len(history.Messages) - 1; i >= 0; i-- {
		if history.Messages[i].Text != "+" && history.Messages[i].Text != "-" {
			lastMessage = history.Messages[i]
			break
		}
	}

	receiver, err := k.users.Find(db, lastMessage.From.ID)
	if err != nil {
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(sendMsg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}
	//// You may not vote on your own message.
	//if update.Message.From.ID == lastMessage.From.ID {
	//	err = errors.New("voting on own message not allowed")
	//	sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	//	_, err := bot.Send(sendMsg)
	//	if err != nil {
	//		log.Printf("error sending message %s\n", err)
	//	}
	//	return
	//}

	kudo, isUpdate, err := k.kudos.UpsertKudo(lastMessage.Text, lastMessage.MessageID, receiver.ID, lastMessage.Chat.ID, k.kudos.IsPositive(update.Message.Text), db)

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	_, err = k.kudoCounts.UpdateKudoCount(kudo, receiver, db, update.Message.Chat.ID, isUpdate)

	if err != nil {
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	if err != nil {
		if _, err := bot.Send(sendMsg); err != nil {
			log.Printf("error sending message %s\n", err)
		}
	}
}

func (k AddKudo) Trigger(update tgbotapi.Update) bool {
	return update.Message.ReplyToMessage == nil && !update.Message.From.IsBot && (strings.EqualFold(update.Message.Text, "+") || strings.EqualFold(update.Message.Text, "-"))
}
