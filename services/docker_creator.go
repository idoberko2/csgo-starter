package services

import (
	"csgo-starter/types"

	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// DockerCreator generates Docker API clients
type DockerCreator struct{}

// Create creates a new Docker API client
func (d *DockerCreator) Create(ip string) (types.Docker, error) {
	dockerClient, err := client.NewClient("http://"+ip+":2375", "1.40", nil, nil)
	if err != nil {
		return nil, err
	}
	log.Debug("Created docker client successfully")

	return &Docker{
		client:   dockerClient,
		progChan: make(chan int, 1),
	}, nil
}
