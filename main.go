package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"niverobot/features"
	"niverobot/model"
	"os"
	"strconv"
)

func getEnvString(name string, value string) string {
	if envValue, ok := os.LookupEnv(name); ok {
		value = envValue
	}
	return value
}
func getEnvInt(name string, value int) int {
	if envValue, ok := os.LookupEnv(name); ok {
		if intVal, err := strconv.Atoi(envValue); err != nil {
			value = intVal
		}
	}
	return value
}

var botToken = getEnvString("BOT_API_TOKEN", "secret")

var dbHost = getEnvString("DB_HOST", "localhost")
var dbPort = getEnvInt("DB_PORT", 5432)
var dbUser = getEnvString("DB_USER", "postgres")
var dbPassword = getEnvString("DB_PASSWORD", "postgres")
var dbSchema = getEnvString("DB_SCHEMA", "postgres")

func main() {
	// TODO: move init db / cleanup
	log.Print("Connecting to database...")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbSchema,
		dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect to database: %s", err)
	}

	log.Println("Applying migrations...")
	err = db.AutoMigrate(&model.User{}, &model.Kudo{}, &model.KudoCount{})
	if err != nil {
		log.Panicf("migration failed: %s", err)
	}

	log.Println("Connecting to Telegram...")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panicf("Failed connecting to Telegram: %s", err)
	}

	// TODO: check this flag
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	kudoService := model.Kudos{}
	userService := model.Users{}
	kudoCountService := model.KudoCounts{}

	var messageHistory = map[int64]model.MessageHistory{}

	// Different features of the bot
	actions := []features.Action{
		features.NewAddKudo(kudoService, userService, kudoCountService),
		features.NewGetKudo(kudoCountService),
		features.NewGetKudoCountOverview(kudoCountService),
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		// Keep a history of messages because useful.
		messageHistory[update.Message.Chat.ID] = messageHistory[update.Message.Chat.ID].AddMessage(update)

		for _, i := range actions {
			// TODO: wrap this in a transaction
			if i.Trigger(update) {
				i.Execute(update, db, bot)
			}
		}
	}
}
