package features

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"niverobot/model"
	"sort"
	"strings"
)

type GetKudoCountOverview struct {
	kudoCounts model.KudoCounts
}

func NewGetKudoCountOverview(kudoCountService model.KudoCounts) GetKudoCountOverview {
	return GetKudoCountOverview{kudoCounts: kudoCountService}
}

func (g GetKudoCountOverview) Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI, history model.MessageHistory) {
	kudoCounts, err := g.kudoCounts.GetKudoCountPerChat(update.Message.Chat.ID, db)
	kudoCounts = orderByKudoSum(kudoCounts)
	overviewMessage := "Current kudo count: \n"
	for _, kudo := range kudoCounts {
		overviewMessage += fmt.Sprintf("%s: %d \n", kudo.User.Username.String, kudo.Plus-kudo.Minus)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s", overviewMessage))
	msg.ParseMode = "html"
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

func orderByKudoSum(kudoCounts []model.KudoCount) []model.KudoCount {
	sort.Slice(kudoCounts, func(i, j int) bool {
		return (kudoCounts[i].Plus - kudoCounts[i].Minus) > (kudoCounts[j].Plus - kudoCounts[j].Minus)
	})
	return kudoCounts
}
