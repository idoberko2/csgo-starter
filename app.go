package main

import (
	"context"
	"csgo-starter/server"
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func reply(ctx context.Context, runner *server.Runner, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var msg tgbotapi.MessageConfig

	log.WithFields(log.Fields{
		"chat": update.Message.Chat,
		"text": update.Message.Text,
	}).Debug("Received message")

	if isValidChat(update.Message.Chat) {
		if update.Message.Text == "/startserver" {
			state, err := runner.Start(ctx)
			if err != nil && errors.As(err, &server.ErrServerStarted{}) {
				errIP := err.(server.ErrServerStarted)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת כבר רץ\n"+errIP.IP)
			} else if err != nil {
				log.WithError(err).Error("An error occurred")
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "קרתה שגיאה")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת מתחיל...\n"+state.DropletIP)
			}
		} else if update.Message.Text == "/stopserver" {
			err := runner.Stop(ctx)
			if err != nil && errors.As(err, &server.ErrServerIdle{}) {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "השרת לא רץ. מה אתה רוצה שאעצור?!")
			} else if err != nil {
				log.WithError(err).Error("An error occurred")
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "קרתה שגיאה")
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "לילה טוב!")
			}
		} else {
			log.WithField("msg", update.Message.Text).Debug("Ignoring")
		}
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Not allowed")
	}

	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

// Run runs the app
func Run() {
	godotenv.Load()
	ctx := context.Background()
	log.SetLevel(log.DebugLevel)

	runner := server.NewRunner()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.WithError(err).Fatal("Error creating Bot API")
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go reply(ctx, runner, bot, update)
	}
}

func isValidChat(chat *tgbotapi.Chat) bool {
	return true
}