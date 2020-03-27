package mocks

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
)

// Sender mock
type Sender struct {
	mock.Mock
}

// Send mock
func (s *Sender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := s.Called(c)

	var arg0 tgbotapi.Message

	if args.Get(0) != nil {
		arg0 = args.Get(0).(tgbotapi.Message)
	}

	return arg0, args.Error(1)
}
