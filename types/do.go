package types

import "context"

// DigitalOcean represents entities that interact with DigitalOcean's API
type DigitalOcean interface {
	StartDroplet(context.Context) (string, int, bool, error)
	StopDroplet(context.Context, int) error
}
