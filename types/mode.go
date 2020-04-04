package types

// Mode represents the current server operation mode
type Mode int

const (
	// ModeIdle server is not running
	ModeIdle Mode = iota

	// ModeStartingDroplet droplet is starting
	ModeStartingDroplet

	// ModeStartedDroplet droplet is started
	ModeStartedDroplet

	// ModeStartingContainer container is starting
	ModeStartingContainer

	// ModeContainerProgress container has started and making progress
	ModeContainerProgress

	// ModeReady server is running and can be joined
	ModeReady

	// ModeShuttingDown server is currently being shut down
	ModeShuttingDown
)
