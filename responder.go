package main

import (
	"context"
	"csgo-starter/types"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Responder is in charge of responding bot messages
type Responder struct {
	Runner types.ServerRunner
	Bot    types.Sender
}

// Respond responds to messages
func (r *Responder) Respond(ctx context.Context, update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig

	log.WithFields(log.Fields{
		"chat": update.Message.Chat,
		"text": update.Message.Text,
	}).Debug("Received message")

	if isValidChat(update.Message.Chat) {
		if update.Message.Text == "/startserver" {
			// initial response since it might be a long action
			r.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "סרג'יו קוסטנזה לשירותך המפקד!"))

			state, err := r.Runner.Start(ctx)
			if err != nil && errors.As(err, &types.ErrServerStarted{}) {
				// server already started
				errIP := err.(types.ErrServerStarted)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת כבר רץ\n"+errIP.IP)
			} else if err != nil {
				// unknown error
				log.WithError(err).Error("An error occurred")
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "קרתה שגיאה")
			} else {
				// starting server
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת מתחיל...\n"+state.DropletIP)
			}
		} else if update.Message.Text == "/stopserver" {
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
		} else {
			log.WithField("msg", update.Message.Text).Debug("Ignoring")
		}
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Not allowed")
	}

	msg.ReplyToMessageID = update.Message.MessageID
	r.Bot.Send(msg)
}

func isValidChat(chat *tgbotapi.Chat) bool {
	return true
}
