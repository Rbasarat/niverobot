package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MessageHistory struct {
	Message []*tgbotapi.Message
}

func (m MessageHistory) AddMessage(update tgbotapi.Update) MessageHistory {
	if len(m.Message) == 2 {
		m.DeleteMessage()
	}
	m.Message = append(m.Message, update.Message)
	return m
}

func (m *MessageHistory) DeleteMessage() *MessageHistory {
	m.Message = m.Message[1:]
	return m
}
