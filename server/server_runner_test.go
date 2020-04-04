package server

import (
	"context"
	"csgo-starter/mocks"
	"csgo-starter/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStart(t *testing.T) {
	ctx := context.Background()
	initState := &types.State{
		DropletID:   0,
		ContainerID: "",
		DropletIP:   "",
		Mode:        types.ModeStartingDroplet,
	}
	dropletCreatedState := &types.State{
		Mode:      types.ModeStartedDroplet,
		DropletID: 1234,
		DropletIP: "2.2.2.2",
	}
	contStartState := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Mode:        types.ModeStartingContainer,
	}
	prog50 := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Progress:    50,
		Mode:        types.ModeContainerProgress,
	}
	prog80 := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Progress:    80,
		Mode:        types.ModeContainerProgress,
	}
	readyState := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Progress:    100,
		Mode:        types.ModeReady,
	}

	stateDAO := mocks.StateDAO{}
	stateDAO.On("SetStartingState").Return(initState, nil)
	stateDAO.On("SetState", dropletCreatedState).Return(dropletCreatedState, nil)
	stateDAO.On("SetState", contStartState).Return(contStartState, nil)
	stateDAO.On("SetState", prog50).Return(prog50, nil)
	stateDAO.On("SetState", prog80).Return(prog80, nil)
	stateDAO.On("SetState", readyState).Return(readyState, nil)

	do := mocks.DigitalOcean{}
	do.On("StartDroplet", ctx).Return("2.2.2.2", 1234, nil)

	docker := mocks.Docker{}
	docker.On("StartContainer", ctx).Return("containerid", nil)
	docker.On("WaitProgress", ctx, mock.Anything).Return(nil)

	dockerCreator := mocks.DockerCreator{}
	dockerCreator.On("Create", "2.2.2.2").Return(&docker, nil)

	runner := Runner{
		stateDAO:      &stateDAO,
		do:            &do,
		dockerCreator: &dockerCreator,
	}

	stateChan := make(chan types.State, 1)
	errChan := make(chan error, 1)

	receivedStates := []types.State{}
	go runner.Start(ctx, stateChan, errChan)

	for stateChan != nil && errChan != nil {
		select {
		case state := <-stateChan:
			{
				receivedStates = append(receivedStates, state)
				if state.Mode == types.ModeReady {
					stateChan = nil
				}
			}
		case err := <-errChan:
			{
				t.Fatal("Unexpected error: ", err)
				errChan = nil
			}
		}
	}

	assert.Equal(t, []types.State{
		*initState,
		*dropletCreatedState,
		*prog50,
		*prog80,
		*readyState,
	}, receivedStates)
	stateDAO.AssertCalled(t, "SetState", contStartState)
}
