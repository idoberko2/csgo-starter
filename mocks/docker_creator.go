package mocks

import (
	"csgo-starter/types"

	"github.com/stretchr/testify/mock"
)

// DockerCreator mock
type DockerCreator struct {
	mock.Mock
}

// Create mock
func (d *DockerCreator) Create(ip string) (types.Docker, error) {
	args := d.Called(ip)

	var arg0 types.Docker

	if args.Get(0) != nil {
		arg0 = args.Get(0).(types.Docker)
	}

	return arg0, args.Error(1)
}
