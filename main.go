package main

import (
	"context"
	"csgo-starter/services"
	"io"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		dockerClient *client.Client
		err          error
	)
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

	dockerClient, err = client.NewClient("http://"+ip+":2375", "1.40", nil, nil)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to docker engine")
	}
	log.Debug("Created docker client successfully")

	err = backoff.Retry(func() error {
		var err error

		reader, err := dockerClient.ImagePull(ctx, "docker.io/cm2network/csgo:latest", types.ImagePullOptions{})
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, reader)
		log.Debug("Pulled image successfully")
		return nil
	}, backoff.NewExponentialBackOff())
	if err != nil {
		log.WithError(err).Fatal("Error pulling image")
	}

	hostconfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"27015": []nat.PortBinding{{HostPort: "27015"}},
			"27020": []nat.PortBinding{{HostPort: "27020"}},
		},
		NetworkMode: container.NetworkMode("host"),
	}

	env := []string{
		"SRCDS_RCONPW=" + os.Getenv("SRCDS_RCONPW"),
		"SRCDS_PW=" + os.Getenv("SRCDS_PW"),
		"SRCDS_PORT=27015",
		"SRCDS_TV_PORT=27020",
		"SRCDS_FPSMAX=300",
		"SRCDS_TICKRATE=128",
		"SRCDS_MAXPLAYERS=14",
		"SRCDS_STARTMAP=de_dust2",
		"SRCDS_REGION=3",
		"SRCDS_MAPGROUP=mg_active",
		"SRCDS_TOKEN=" + os.Getenv("SRCDS_TOKEN"),
	}

	log.WithField("env", env).Info("Starting container...")

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "cm2network/csgo",
		Env:   env,
	}, hostconfig, nil, "cs")
	if err != nil {
		log.WithError(err).Fatal("Error creating container")
	}

	log.Debug("Created docker container successfully")

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.WithError(err).Fatal("Error starting container")
	}

	log.Debug("Started docker container successfully")
}
