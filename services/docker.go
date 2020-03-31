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

//  Update state (0x61) downloading, progress: 98.47 (24447265031 / 24827773544)
//  Update state (0x61) downloading, progress: 98.65 (24492240398 / 24827773544)
//  Update state (0x61) downloading, progress: 98.81 (24532047166 / 24827773544)
//  Update state (0x61) downloading, progress: 98.96 (24569868334 / 24827773544)
//  Update state (0x61) downloading, progress: 99.07 (24597954452 / 24827773544)
//  Update state (0x61) downloading, progress: 99.19 (24626266004 / 24827773544)
//  Update state (0x61) downloading, progress: 99.40 (24677958411 / 24827773544)
//  Update state (0x61) downloading, progress: 99.56 (24718313668 / 24827773544)
//  Update state (0x61) downloading, progress: 99.68 (24747560956 / 24827773544)
//  Update state (0x61) downloading, progress: 99.87 (24796179593 / 24827773544)
// RecordSteamInterfaceCreation (PID 67): SteamGameServer012 / GameServer
// RecordSteamInterfaceCreation (PID 67): SteamUtils008 / Utils
// RecordSteamInterfaceCreation (PID 67): SteamNetworking005 / Networking
// RecordSteamInterfaceCreation (PID 67): SteamGameServerStats001 / GameServerStats
// RecordSteamInterfaceCreation (PID 67): STEAMHTTP_INTERFACE_VERSION002 / HTTP
// RecordSteamInterfaceCreation (PID 67): STEAMINVENTORY_INTERFACE_V001 / Inventory
// RecordSteamInterfaceCreation (PID 67): STEAMUGC_INTERFACE_VERSION008 / UGC
// RecordSteamInterfaceCreation (PID 67): STEAMAPPS_INTERFACE_VERSION008 / Apps
// RecordSteamInterfaceCreation (PID 67): SteamUtils009 / Utils
// RecordSteamInterfaceCreation (PID 67): SteamNetworkingSocketsSerialized003 /
// RecordSteamInterfaceCreation (PID 67): SteamGameServer012 / GameServer
// RecordSteamInterfaceCreation (PID 67): STEAMHTTP_INTERFACE_VERSION003 / HTTP
// RecordSteamInterfaceCreation (PID 67): SteamGameServer012 / GameServer
// RecordSteamInterfaceCreation (PID 67): SteamUtils008 / Utils
// RecordSteamInterfaceCreation (PID 67): SteamNetworking005 / Networking
// RecordSteamInterfaceCreation (PID 67): SteamGameServerStats001 / GameServerStats
// RecordSteamInterfaceCreation (PID 67): STEAMHTTP_INTERFACE_VERSION002 / HTTP
// RecordSteamInterfaceCreation (PID 67): STEAMINVENTORY_INTERFACE_V001 / Inventory
// RecordSteamInterfaceCreation (PID 67): STEAMUGC_INTERFACE_VERSION008 / UGC
// RecordSteamInterfaceCreation (PID 67): STEAMAPPS_INTERFACE_VERSION008 / Apps
// RecordSteamInterfaceCreation (PID 67): SteamGameCoordinator001 /
// RecordSteamInterfaceCreation (PID 67): SteamGameServer012 / GameServer

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
