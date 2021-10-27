package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MessageHistory struct {
	Messages []*tgbotapi.Message
}

func (m MessageHistory) AddMessage(update tgbotapi.Update) MessageHistory {
	// TODO: change this to config
	if len(m.Messages) == 10 {
		m.DeleteMessage()
	}
	m.Messages = append(m.Messages, update.Message)
	return m
}

func (m *MessageHistory) DeleteMessage() *MessageHistory {
	m.Messages = m.Messages[1:]
	return m
}


