package features

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

var lastUpdate = make(map[int64]tgbotapi.Update)

func NewAddKudo(kudosService model.Kudos, userService model.Users, kudoCountService model.KudoCounts) AddKudo {

	return AddKudo{kudos: kudosService, users: userService, kudoCounts: kudoCountService}
}

func (k AddKudo) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI) {
	if lastUpdate[update.Message.Chat.ID].Message == nil {
		return
	}
	var msg tgbotapi.MessageConfig
	_, err := k.users.CreateUserIfNotExist(update.Message.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	receiver, err := k.users.CreateUserIfNotExist(lastUpdate[update.Message.Chat.ID].Message.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))

	}

	if lastUpdate[update.Message.Chat.ID].Message.From.ID == update.Message.From.ID {
		err = errors.New("voting on own message not allowed")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}

	kudo, err := k.kudos.CreateKudoIfNotExist(lastUpdate[update.Message.Chat.ID].Message, update.Message, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}
	_, err = k.kudoCounts.UpdateKudoCount(kudo, receiver, db, update.Message.Chat.ID)

	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	if err != nil {
		if _, err := bot.Send(msg); err != nil {
			log.Printf("error sending message %s\n", err)
		}
	}

}

func (k AddKudo) Trigger(update tgbotapi.Update) bool {
	trigger :=  update.Message.ReplyToMessage == nil && !update.Message.From.IsBot && (strings.EqualFold(update.Message.Text, "+") || strings.EqualFold(update.Message.Text, "-"))
	if !trigger {
		lastUpdate[update.Message.Chat.ID] = update
	}
	return trigger
}
