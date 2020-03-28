package types

import "context"

// Docker represents entities that interact with Docker's API
type Docker interface {
	StartContainer(ctx context.Context) (string, error)
}
