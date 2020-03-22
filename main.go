package main

import (
	"context"
	"csgo-starter/services"
	"os"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	ctx := context.Background()
	log.SetLevel(log.DebugLevel)
	do := services.NewDo(os.Getenv("DO_TOKEN"))
	ip, did, err := do.StartDroplet(ctx)
	if err != nil {
		log.WithError(err).Fatal("Error creating droplet")
	}

	log.WithFields(log.Fields{
		"ip":        ip,
		"dropletID": did,
	}).Info("Started droplet")

	dock, err := services.NewDocker(ip)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to docker engine")
	}

	containerID, err := dock.StartContainer(ctx)
	if err != nil {
		log.WithError(err).Fatal("Error initializing docker container")
	}

	_ = containerID
}
