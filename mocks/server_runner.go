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
func (r *ServerRunner) Start(ctx context.Context, stateChan chan types.State, errChan chan error) {
	r.Called(ctx, stateChan, errChan)
}

// Stop mock
func (r *ServerRunner) Stop(ctx context.Context) error {
	return r.Called(ctx).Error(0)
}
