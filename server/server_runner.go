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
	stateDAO *StateDAO
}

// NewRunner generates a new instance of Runner
func NewRunner() *Runner {
	return &Runner{
		stateDAO: NewStateDAO(),
	}
}

// Start starts the server
func (rnr *Runner) Start(ctx context.Context) (*types.State, error) {
	_, err := rnr.stateDAO.SetStartingState()
	if err != nil {
		return nil, err
	}

	do := services.NewDo(os.Getenv("DO_TOKEN"))
	ip, did, err := do.StartDroplet(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating droplet")
	}

	log.WithFields(log.Fields{
		"ip":        ip,
		"dropletID": did,
	}).Info("Started droplet")
	_, err = rnr.stateDAO.SetState(&types.State{
		Mode:      types.ModeStarting,
		DropletID: did,
		DropletIP: ip,
	})
	if err != nil {
		return nil, err
	}

	dock, err := services.NewDocker(ip)
	if err != nil {
		return nil, errors.Wrap(err, "Error connecting to docker engine")
	}

	containerID, err := dock.StartContainer(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error initializing docker container")
	}

	log.WithFields(log.Fields{
		"containerID": containerID,
	}).Info("Started container")

	return rnr.stateDAO.SetState(&types.State{
		Mode:        types.ModeStarted,
		DropletID:   did,
		DropletIP:   ip,
		ContainerID: containerID,
	})
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

	do := services.NewDo(os.Getenv("DO_TOKEN"))
	err = do.StopDroplet(ctx, state.DropletID)
	if err != nil {
		return err
	}

	_, err = rnr.stateDAO.SetState(&types.State{})
	if err != nil {
		return err
	}

	return nil
}
