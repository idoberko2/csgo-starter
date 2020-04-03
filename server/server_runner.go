package server

import (
	"context"
	"csgo-starter/services"
	"csgo-starter/types"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Runner is in charge of running the CS:GO server
type Runner struct {
	stateDAO      types.StateDAO
	dockerCreator types.DockerCreator
	do            types.DigitalOcean
}

// NewRunner generates a new instance of Runner
func NewRunner() *Runner {
	return &Runner{
		stateDAO:      NewStateDAO(),
		dockerCreator: &services.DockerCreator{},
		do:            services.NewDo(os.Getenv("DO_TOKEN")),
	}
}

// Start starts the server
func (rnr *Runner) Start(ctx context.Context, stateChan chan types.State, errChan chan error) {
	state, err := rnr.stateDAO.SetStartingState()
	if err != nil {
		errChan <- err
		return
	}
	stateChan <- *state

	ip, did, err := rnr.do.StartDroplet(ctx)
	if err != nil {
		errChan <- errors.Wrap(err, "Error creating droplet")
		return
	}

	log.WithFields(log.Fields{
		"ip":        ip,
		"dropletID": did,
	}).Info("Started droplet")
	state, err = rnr.stateDAO.SetState(&types.State{
		Mode:      types.ModeStartedDroplet,
		DropletID: did,
		DropletIP: ip,
	})
	if err != nil {
		errChan <- err
		return
	}
	stateChan <- *state

	dock, err := rnr.dockerCreator.Create(ip)
	if err != nil {
		errChan <- errors.Wrap(err, "Error connecting to docker engine")
		return
	}

	containerID, err := dock.StartContainer(ctx)
	if err != nil {
		errChan <- errors.Wrap(err, "Error initializing docker container")
		return
	}

	log.WithFields(log.Fields{
		"containerID": containerID,
	}).Info("Started container")

	state, err = rnr.stateDAO.SetState(&types.State{
		Mode:        types.ModeStartingContainer,
		DropletID:   did,
		DropletIP:   ip,
		ContainerID: containerID,
	})
	if err != nil {
		errChan <- err
		return
	}

	stateChan <- *state
}

// Stop stops the CS:GO server
func (rnr *Runner) Stop(ctx context.Context) error {
	state, err := rnr.stateDAO.GetState()
	if err != nil {
		return err
	}

	if state.Mode == types.ModeIdle {
		return types.ErrServerIdle{}
	}

	err = rnr.do.StopDroplet(ctx, state.DropletID)
	if err != nil {
		return err
	}

	_, err = rnr.stateDAO.SetState(&types.State{})
	if err != nil {
		return err
	}

	return nil
}
