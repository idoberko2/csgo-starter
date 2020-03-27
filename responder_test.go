package main

import (
	"context"
	"csgo-starter/mocks"
	"csgo-starter/types"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
)

func TestRespond_start(t *testing.T) {
	bot := mocks.Sender{}
	bot.On("Send", mock.Anything).Return(nil, nil)

	chatid := int64(1111)
	messageid := 2222

	testCases := []struct {
		desc          string
		returnedError error
		expected      []tgbotapi.MessageConfig
	}{
		{
			desc: "start - already running",
			returnedError: types.ErrServerStarted{
				IP: "1.1.1.1",
			},
			expected: []tgbotapi.MessageConfig{
				tgbotapi.NewMessage(chatid, "סרג'יו קוסטנזה לשירותך המפקד!"),
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "השרת כבר רץ\n" + "1.1.1.1",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			runner := mocks.ServerRunner{}
			runner.On("Start", mock.Anything).Return(nil, tC.returnedError)

			responder := Responder{
				Runner: &runner,
				Bot:    &bot,
			}

			responder.Respond(context.Background(), tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text:      "/startserver",
					MessageID: messageid,
					Chat: &tgbotapi.Chat{
						ID: chatid,
					},
				},
			})

			runner.AssertCalled(t, "Start", context.Background())
			bot.AssertNumberOfCalls(t, "Send", len(tC.expected))

			for _, exMsg := range tC.expected {
				bot.AssertCalled(t, "Send", exMsg)
			}
		})
	}
}
