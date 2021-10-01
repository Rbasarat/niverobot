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

func NewAddKudo(kudosService model.Kudos, userService model.Users, kudoCountService model.KudoCounts) AddKudo {
	return AddKudo{kudos: kudosService, users: userService, kudoCounts: kudoCountService}
}

func (k AddKudo) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	receiver, err := k.users.CreateUserIfNotExist(update.Message.ReplyToMessage.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	_, err = k.users.CreateUserIfNotExist(update.Message.From, db)

	// You may not vote on your own message.
	if strings.EqualFold(update.Message.ReplyToMessage.Text, "+") || strings.EqualFold(update.Message.ReplyToMessage.Text, "-") {
		err = errors.New("voting on a kudo not allowed")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}
	if update.Message.From.ID == update.Message.ReplyToMessage.From.ID {
		err = errors.New("voting on own message not allowed")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	kudo, err := k.kudos.CreateFromReplyKudoIfNotExist(update.Message, db)
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
	return update.Message.ReplyToMessage != nil && !update.Message.ReplyToMessage.From.IsBot && (strings.EqualFold(update.Message.Text, "+") || strings.EqualFold(update.Message.Text, "-"))
}
