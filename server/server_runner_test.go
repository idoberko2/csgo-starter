package server

import (
	"context"
	"csgo-starter/mocks"
	"csgo-starter/types"
	"testing"

	"github.com/stretchr/testify/assert"
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
	exState := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Mode:        types.ModeStartingContainer,
	}

	stateDAO := mocks.StateDAO{}
	stateDAO.On("SetStartingState").Return(initState, nil)
	stateDAO.On("SetState", dropletCreatedState).Return(dropletCreatedState, nil)
	stateDAO.On("SetState", exState).Return(exState, nil)

	do := mocks.DigitalOcean{}
	do.On("StartDroplet", ctx).Return("2.2.2.2", 1234, nil)

	docker := mocks.Docker{}
	docker.On("StartContainer", ctx).Return("containerid", nil)

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
				t.Log(state)
				receivedStates = append(receivedStates, state)
				if state.Mode == types.ModeStartingContainer {
					stateChan = nil
					// t.FailNow()
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
		*exState,
	}, receivedStates)
	stateDAO.AssertCalled(t, "SetState", exState)
}
