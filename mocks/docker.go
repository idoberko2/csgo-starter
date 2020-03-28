package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Docker mock
type Docker struct {
	mock.Mock
}

// StartContainer mock
func (d *Docker) StartContainer(ctx context.Context) (string, error) {
	args := d.Called(ctx)

	return args.String(0), args.Error(1)
}
