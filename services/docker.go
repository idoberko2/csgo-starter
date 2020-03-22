package services

import (
	"context"
	"io"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Docker is in charge of interacting with docker engine API
type Docker struct {
	client *client.Client
}

// NewDocker generates a new Docker instance
func NewDocker(ip string) (*Docker, error) {
	dockerClient, err := client.NewClient("http://"+ip+":2375", "1.40", nil, nil)
	if err != nil {
		return nil, err
	}
	log.Debug("Created docker client successfully")

	return &Docker{
		client: dockerClient,
	}, nil
}

// StartContainer starts the CS:GO container
func (dock *Docker) StartContainer(ctx context.Context) (string, error) {
	err := dock.waitAndPull(ctx)
	if err != nil {
		return "", errors.Wrap(err, "Error pulling image")
	}

	resp, err := dock.client.ContainerCreate(ctx, &container.Config{
		Image: "cm2network/csgo",
		Env:   getEnv(),
	}, getHostConfig(), nil, "cs")
	if err != nil {
		return "", errors.Wrap(err, "Error creating container")
	}

	log.Debug("Created docker container successfully")

	if err := dock.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.Wrap(err, "Error starting container")
	}

	log.Debug("Started docker container successfully")

	return resp.ID, nil
}

func (dock *Docker) waitAndPull(ctx context.Context) error {
	return backoff.Retry(func() error {
		var err error

		reader, err := dock.client.ImagePull(ctx, "docker.io/cm2network/csgo:latest", types.ImagePullOptions{})
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, reader)
		log.Debug("Pulled image successfully")
		return nil
	}, backoff.NewExponentialBackOff())
}

func getEnv() []string {
	return []string{
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
}

func getHostConfig() *container.HostConfig {
	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"27015": []nat.PortBinding{{HostPort: "27015"}},
			"27020": []nat.PortBinding{{HostPort: "27020"}},
		},
		NetworkMode: container.NetworkMode("host"),
	}
}
