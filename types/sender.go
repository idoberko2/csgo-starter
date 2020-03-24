package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Sender represents entities that are able to send messages
type Sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
