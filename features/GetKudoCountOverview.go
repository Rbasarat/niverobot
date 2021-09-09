package features

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olekukonko/tablewriter"
	"gorm.io/gorm"
	"log"
	"niverobot/model"
	"sort"
	"strconv"
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
	kudoCounts = orderByKudoSum(kudoCounts)
	table := fmt.Sprintf("%s", renderTableAsString(kudoCounts))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("<pre>%s</pre>", table))
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

func renderTableAsString(kudoCount []model.KudoCount) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"Name", "Plus", "Min"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, kudo := range kudoCount {
		table.Append([]string{kudo.User.FirstName.String, strconv.Itoa(kudo.Plus), strconv.Itoa(kudo.Minus)})
	}
	table.Render()
	return tableString.String()
}
