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

type AddKudoFromReply struct {
	kudos      model.Kudos
	users      model.Users
	kudoCounts model.KudoCounts
}

func NewAddKudoFromReply(kudosService model.Kudos, userService model.Users, kudoCountService model.KudoCounts) AddKudoFromReply {
	return AddKudoFromReply{kudos: kudosService, users: userService, kudoCounts: kudoCountService}
}

func (k AddKudoFromReply) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI, history model.MessageHistory) {
	var sendMsg tgbotapi.MessageConfig

	receiver, err := k.users.Find(db, update.Message.ReplyToMessage.From.ID)
	if err != nil {
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(sendMsg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
	}

	if strings.EqualFold(update.Message.ReplyToMessage.Text, "+") || strings.EqualFold(update.Message.ReplyToMessage.Text, "-") {
		err = errors.New("voting on a kudo not allowed")
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(sendMsg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}

	//You may not vote on your own message.
	if update.Message.From.ID == update.Message.ReplyToMessage.From.ID {
		err = errors.New("voting on own message not allowed")
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(sendMsg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
		return
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		sendMsg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
		_, err = bot.Send(sendMsg)
		if err != nil {
			log.Printf("error sending message %s\n", err)
		}
	}

	kudo, isUpdate, err := k.kudos.UpsertKudo(update.Message.ReplyToMessage.Text, update.Message.ReplyToMessage.MessageID, update.Message.ReplyToMessage.From.ID, update.Message.Chat.ID, k.kudos.IsPositive(update.Message.Text), db)

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

func (k AddKudoFromReply) Trigger(update tgbotapi.Update) bool {
	return update.Message.ReplyToMessage != nil && !update.Message.ReplyToMessage.From.IsBot && (strings.EqualFold(update.Message.Text, "+") || strings.EqualFold(update.Message.Text, "-"))
}
