package types

import "context"

// ServerRunner represents entities that run a server
type ServerRunner interface {
	Start(ctx context.Context) (*State, error)
	Stop(ctx context.Context) error
}
