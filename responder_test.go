package main

import (
	"context"
	"csgo-starter/mocks"
	"csgo-starter/types"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
)

func TestRespond_start(t *testing.T) {
	chatid := int64(1111)
	messageid := 2222

	testCases := []struct {
		desc          string
		returnedState *types.State
		returnedError error
		expected      []tgbotapi.MessageConfig
	}{
		{
			desc:          "start - already running",
			returnedState: nil,
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
		{
			desc:          "start - error",
			returnedState: nil,
			returnedError: errors.New("some network"),
			expected: []tgbotapi.MessageConfig{
				tgbotapi.NewMessage(chatid, "סרג'יו קוסטנזה לשירותך המפקד!"),
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "קרתה שגיאה",
				},
			},
		},
		{
			desc: "start - ok",
			returnedState: &types.State{
				DropletIP: "1.1.1.1",
			},
			returnedError: nil,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.NewMessage(chatid, "סרג'יו קוסטנזה לשירותך המפקד!"),
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "השרת מתחיל...\n1.1.1.1",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bot := mocks.Sender{}
			bot.On("Send", mock.Anything).Return(nil, nil)

			runner := mocks.ServerRunner{}
			runner.On("Start", mock.Anything).Return(tC.returnedState, tC.returnedError)

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
