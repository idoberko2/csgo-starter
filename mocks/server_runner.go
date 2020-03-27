package mocks

import (
	"context"
	"csgo-starter/types"

	"github.com/stretchr/testify/mock"
)

// ServerRunner mock
type ServerRunner struct {
	mock.Mock
}

// Start mock
func (r *ServerRunner) Start(ctx context.Context) (*types.State, error) {
	args := r.Called(ctx)

	var arg0 *types.State

	if args.Get(0) != nil {
		arg0 = args.Get(0).(*types.State)
	}

	return arg0, args.Error(1)
}

// Stop mock
func (r *ServerRunner) Stop(ctx context.Context) error {
	return r.Called(ctx).Error(0)
}
