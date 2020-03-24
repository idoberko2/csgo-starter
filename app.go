package main

import (
	"context"
	"csgo-starter/server"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

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

	responder := &Responder{
		Runner: runner,
		Bot:    bot,
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go responder.Respond(ctx, update)
	}
}
