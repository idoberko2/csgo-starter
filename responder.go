package main

import (
	"context"
	"csgo-starter/types"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Responder is in charge of responding bot messages
type Responder struct {
	Runner types.ServerRunner
	Bot    types.Sender
}

func (r *Responder) handleStart(ctx context.Context, update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig

	stateChan := make(chan types.State, 1)
	errChan := make(chan error, 1)

	// initial response since it might be a long action
	r.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "סרג'יו קוסטנזה לשירותך המפקד!"))

	go r.Runner.Start(ctx, stateChan, errChan)

	for stateChan != nil && errChan != nil {
		select {
		case state := <-stateChan:
			{
				if state.Mode == types.ModeStartedDroplet {
					// starting server
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת מתחיל...\n"+state.DropletIP)
					stateChan = nil
				}
			}
		case err := <-errChan:
			{
				if errors.As(err, &types.ErrServerStarted{}) {
					// server already started
					errIP := err.(types.ErrServerStarted)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת כבר רץ\n"+errIP.IP)
				} else {
					// unknown error
					log.WithError(err).Error("An error occurred")
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "קרתה שגיאה")
				}

				errChan = nil
			}
		}
	}

	msg.ReplyToMessageID = update.Message.MessageID
	r.Bot.Send(msg)
}

func (r *Responder) handleStop(ctx context.Context, update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig
	err := r.Runner.Stop(ctx)
	if err != nil && errors.As(err, &types.ErrServerIdle{}) {
		// server is not running
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת לא רץ. מה אתה רוצה שאעצור?!")
	} else if err != nil {
		// unknown error
		log.WithError(err).Error("An error occurred")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "קרתה שגיאה")
	} else {
		// stopping server
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "לילה טוב!")
	}

	msg.ReplyToMessageID = update.Message.MessageID
	r.Bot.Send(msg)
}

// Respond responds to messages
func (r *Responder) Respond(ctx context.Context, update tgbotapi.Update) {
	log.WithFields(log.Fields{
		"chat": update.Message.Chat,
		"text": update.Message.Text,
	}).Debug("Received message")

	if isValidChat(update.Message.Chat) {
		if update.Message.Text == "/startserver" || update.Message.Text == "/startserver@"+os.Getenv("TG_BOT_NAME") {
			r.handleStart(ctx, update)
			return
		} else if update.Message.Text == "/stopserver" || update.Message.Text == "/stopserver@"+os.Getenv("TG_BOT_NAME") {
			r.handleStop(ctx, update)
			return
		} else {
			log.WithField("msg", update.Message.Text).Debug("Ignoring")
			return
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Not allowed")
	msg.ReplyToMessageID = update.Message.MessageID
	r.Bot.Send(msg)
}

func isValidChat(chat *tgbotapi.Chat) bool {
	allowedChatIDs := os.Getenv("ALLOWED_CHAT_IDS")
	if allowedChatIDs == "" {
		return true
	}

	for _, strchatid := range strings.Split(allowedChatIDs, ",") {
		chatid, err := strconv.ParseInt(strchatid, 10, 64)
		if err != nil {
			log.WithField("chatid", strchatid).Error("Invalid chat id")
			return false
		}

		if chat.ID == chatid {
			return true
		}
	}

	log.WithField("chatid", chat.ID).Debug("Chat ID is not allowed")

	return false
}
