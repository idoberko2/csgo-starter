package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

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
	client   *client.Client
	progChan chan int
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

	go dock.checkProgress(ctx, resp.ID)

	return resp.ID, nil
}

func (dock *Docker) checkProgress(ctx context.Context, cid string) error {
	reader, err := dock.client.ContainerLogs(ctx, cid, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		log.WithError(err).Error("Error getting docker logs")
		return err
	}
	defer reader.Close()

	progRegex := regexp.MustCompile(`Update state \(0x\d+\) downloading, progress: (?P<progress>\d*\.\d*).*`)
	doneRegex := regexp.MustCompile(`weapon_sound_falloff_multiplier \- 1\.0`)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.WithField("line", scanner.Text()).Debug("Scanned line")
		if progRegex.Match(scanner.Bytes()) {
			matches := progRegex.FindStringSubmatch(scanner.Text())
			log.WithField("matches", matches).Debug("Found matches")
			progStr := matches[1]
			progress, err := strconv.ParseFloat(progStr, 64)
			if err != nil {
				log.WithError(err).Error("Error parsing progress")
				return err
			}
			log.WithField("progress", progress).Debug("Updating progress")
			dock.progChan <- int(progress)
		}

		if doneRegex.Match(scanner.Bytes()) {
			dock.progChan <- 100
			log.Debug("Last progress line")
			break
		}
	}

	return nil
}

// WaitProgress returns when the container starting progress reaches n
func (dock *Docker) WaitProgress(ctx context.Context, n int) error {
	for curProgress := range dock.progChan {
		log.WithField("curProgress", curProgress).Debug("Progress...")
		if curProgress >= n {
			return nil
		}
	}

	return fmt.Errorf("Channel is closed without reaching %d", n)
}

func (dock *Docker) waitAndPull(ctx context.Context) error {
	return backoff.Retry(func() error {
		var err error

		log.Debug("Trying to pull image...")
		reader, err := dock.client.ImagePull(ctx, "docker.io/cm2network/csgo:latest", types.ImagePullOptions{})
		if err != nil {
			log.WithError(err).Debug("Failed pulling image")
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
