package main

import (
	"context"
	"csgo-starter/mocks"
	"csgo-starter/types"
	"errors"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/mock"
)

func TestRespond_start(t *testing.T) {
	chatid := int64(1111)
	messageid := 2222

	testCases := []struct {
		desc          string
		input         string
		returnedState *types.State
		returnedError error
		expected      []tgbotapi.MessageConfig
	}{
		{
			desc:          "start - already running",
			input:         "/startserver",
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
			input:         "/startserver",
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
			desc:  "start - ok",
			input: "/startserver",
			returnedState: &types.State{
				DropletIP: "1.1.1.1",
				Mode:      types.ModeStartedDroplet,
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
		{
			desc:  "start - with mention",
			input: "/startserver@botname",
			returnedState: &types.State{
				DropletIP: "1.1.1.1",
				Mode:      types.ModeStartedDroplet,
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
			os.Setenv("TG_BOT_NAME", "botname")
			bot := mocks.Sender{}
			bot.On("Send", mock.Anything).Return(nil, nil)

			runner := mocks.ServerRunner{}
			runner.On("Start", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				stateChan := args.Get(1).(chan types.State)
				errChan := args.Get(2).(chan error)

				if tC.returnedState != nil {
					stateChan <- *tC.returnedState
				}
				if tC.returnedError != nil {
					errChan <- tC.returnedError
				}
			})

			responder := Responder{
				Runner: &runner,
				Bot:    &bot,
			}

			responder.Respond(context.Background(), tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text:      tC.input,
					MessageID: messageid,
					Chat: &tgbotapi.Chat{
						ID: chatid,
					},
				},
			})

			runner.AssertCalled(t, "Start", context.Background(), mock.Anything, mock.Anything)
			if !bot.AssertNumberOfCalls(t, "Send", len(tC.expected)) {
				t.FailNow()
			}

			for _, exMsg := range tC.expected {
				bot.AssertCalled(t, "Send", exMsg)
			}
		})
	}
}

func TestRespond_stop(t *testing.T) {
	chatid := int64(3333)
	messageid := 4444

	testCases := []struct {
		desc           string
		input          string
		allowedChatIDs string
		returnedError  error
		expected       []tgbotapi.MessageConfig
	}{
		{
			desc:          "stop - not running",
			input:         "/stopserver",
			returnedError: types.ErrServerIdle{},
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "השרת לא רץ. מה אתה רוצה שאעצור?!",
				},
			},
		},
		{
			desc:          "stop - error",
			input:         "/stopserver",
			returnedError: errors.New("some network error"),
			expected: []tgbotapi.MessageConfig{
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
			desc:          "stop - ok",
			input:         "/stopserver",
			returnedError: nil,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "לילה טוב!",
				},
			},
		},
		{
			desc:          "stop - with mention",
			input:         "/stopserver@botname",
			returnedError: nil,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "לילה טוב!",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			os.Setenv("TG_BOT_NAME", "botname")
			bot := mocks.Sender{}
			bot.On("Send", mock.Anything).Return(nil, nil)

			runner := mocks.ServerRunner{}
			runner.On("Stop", mock.Anything).Return(tC.returnedError)

			responder := Responder{
				Runner: &runner,
				Bot:    &bot,
			}

			responder.Respond(context.Background(), tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text:      tC.input,
					MessageID: messageid,
					Chat: &tgbotapi.Chat{
						ID: chatid,
					},
				},
			})

			runner.AssertCalled(t, "Stop", context.Background())
			if !bot.AssertNumberOfCalls(t, "Send", len(tC.expected)) {
				t.FailNow()
			}

			for _, exMsg := range tC.expected {
				bot.AssertCalled(t, "Send", exMsg)
			}
		})
	}
}

func TestRespond_ignore(t *testing.T) {
	chatid := int64(5555)
	messageid := 6666

	bot := mocks.Sender{}
	bot.On("Send", mock.Anything).Return(nil, nil)

	runner := mocks.ServerRunner{}
	runner.On("Start", mock.Anything).Return(nil, nil)
	runner.On("Stop", mock.Anything).Return(nil)

	responder := Responder{
		Runner: &runner,
		Bot:    &bot,
	}

	responder.Respond(context.Background(), tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text:      "some message",
			MessageID: messageid,
			Chat: &tgbotapi.Chat{
				ID: chatid,
			},
		},
	})

	runner.AssertNotCalled(t, "Start", mock.Anything)
	runner.AssertNotCalled(t, "Stop", mock.Anything)
	bot.AssertNotCalled(t, "Send", mock.Anything)
}

func TestRespond_allowed(t *testing.T) {
	chatid := int64(7777)
	messageid := 8888

	testCases := []struct {
		desc           string
		allowedChatIDs string
		isStopCalled   bool
		expected       []tgbotapi.MessageConfig
	}{
		{
			desc:           "allowed",
			allowedChatIDs: "1234,7777",
			isStopCalled:   true,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "לילה טוב!",
				},
			},
		},
		{
			desc:           "not allowed",
			allowedChatIDs: "1234",
			isStopCalled:   false,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "Not allowed",
				},
			},
		},
		{
			desc:           "invalid env variable",
			allowedChatIDs: "should be int64",
			isStopCalled:   false,
			expected: []tgbotapi.MessageConfig{
				tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatid,
						ReplyToMessageID: messageid,
					},
					Text: "Not allowed",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			os.Setenv("ALLOWED_CHAT_IDS", tC.allowedChatIDs)
			bot := mocks.Sender{}
			bot.On("Send", mock.Anything).Return(nil, nil)

			runner := mocks.ServerRunner{}
			runner.On("Stop", mock.Anything).Return(nil)

			responder := Responder{
				Runner: &runner,
				Bot:    &bot,
			}

			responder.Respond(context.Background(), tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text:      "/stopserver",
					MessageID: messageid,
					Chat: &tgbotapi.Chat{
						ID: chatid,
					},
				},
			})

			if tC.isStopCalled {
				runner.AssertCalled(t, "Stop", context.Background())
			} else {
				runner.AssertNotCalled(t, "Stop", mock.Anything)
			}

			if !bot.AssertNumberOfCalls(t, "Send", len(tC.expected)) {
				t.FailNow()
			}

			for _, exMsg := range tC.expected {
				bot.AssertCalled(t, "Send", exMsg)
			}
		})
	}
}
