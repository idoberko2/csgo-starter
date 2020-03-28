package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// DigitalOcean mock
type DigitalOcean struct {
	mock.Mock
}

// StartDroplet mock
func (do *DigitalOcean) StartDroplet(ctx context.Context) (string, int, error) {
	args := do.Called(ctx)

	return args.String(0), args.Int(1), args.Error(2)
}

// StopDroplet mock
func (do *DigitalOcean) StopDroplet(ctx context.Context, dropletID int) error {
	args := do.Called(ctx, dropletID)

	return args.Error(0)
}
