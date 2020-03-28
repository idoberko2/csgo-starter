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
	exState := &types.State{
		DropletID:   1234,
		ContainerID: "containerid",
		DropletIP:   "2.2.2.2",
		Mode:        types.ModeStarted,
	}

	stateDAO := mocks.StateDAO{}
	stateDAO.On("SetStartingState").Return(nil, nil)
	stateDAO.On("SetState", mock.Anything).Return(exState, nil)

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

	actualState, err := runner.Start(ctx)
	assert.Nil(t, err)
	assert.Equal(t, exState, actualState)
	stateDAO.AssertCalled(t, "SetState", exState)
}
