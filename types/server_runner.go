package types

import "context"

// ServerRunner represents entities that run a server
type ServerRunner interface {
	Start(ctx context.Context, stateChan chan State, errChan chan error)
	Stop(ctx context.Context) error
}
