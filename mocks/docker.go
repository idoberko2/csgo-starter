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

// WaitProgress mock
func (d *Docker) WaitProgress(ctx context.Context, n int) error {
	args := d.Called(ctx, n)

	return args.Error(0)
}
