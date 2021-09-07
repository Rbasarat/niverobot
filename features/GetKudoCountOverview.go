package features

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"gorm.io/gorm"
	"log"
	"niverobot/model"
	"strings"
)

type GetKudoCountOverview struct {
	kudoCounts model.KudoCounts
}

func NewGetKudoCountOverview(kudoCountService model.KudoCounts) GetKudoCountOverview {
	return GetKudoCountOverview{kudoCounts: kudoCountService}
}

func (g GetKudoCountOverview) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI) {
	kudoCounts, err := g.kudoCounts.GetKudoCountPerChat(update.Message.Chat.ID, db)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Not implemented yet.")
	msg.ReplyToMessageID = update.Message.MessageID
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Kudo error: %s", err))
	}

	_, err = bot.Send(msg)
	if err != nil {
		log.Panicf("error sending message %s", err)
	}
}

func (g GetKudoCountOverview) Trigger(update tgbotapi.Update) bool {
	return strings.EqualFold(update.Message.Text, ".kudos")
}
