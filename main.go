package main

import (
	"context"
	"csgo-starter/server"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	log.SetLevel(log.DebugLevel)

	runner := server.NewRunner()

	err := runner.Start(ctx)
	if err != nil {
		log.WithError(err).Error("Error starting server")
	} else {
		log.Info("Success")
	}
}
