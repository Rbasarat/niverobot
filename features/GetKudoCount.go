package features

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"niverobot/model"
	"strings"
)

type GetKudoCount struct {
	kudoCounts model.KudoCounts
}

func NewGetKudo(kudoCountService model.KudoCounts) GetKudoCount {
	return GetKudoCount{kudoCounts: kudoCountService}
}

func (g GetKudoCount) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI, history model.MessageHistory) {
	kudoCount, err := g.kudoCounts.GetKudoCount(update.Message.From.ID, update.Message.Chat.ID, db)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Plus kudos: %d\nMin kudos: %d\n", kudoCount.Plus, kudoCount.Minus))
	msg.ReplyToMessageID = update.Message.MessageID
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	_, err = bot.Send(msg)
	if err != nil {
		log.Panicf("error sending message %s", err)
	}
}

func (g GetKudoCount) Trigger(update tgbotapi.Update) bool {
	return strings.EqualFold(update.Message.Text, ".kudo")
}
