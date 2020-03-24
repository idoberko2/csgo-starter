package types

// Mode represents the current server operation mode
type Mode int

const (
	// ModeIdle server is not running
	ModeIdle Mode = iota

	// ModeStarting server is starting
	ModeStarting

	// ModeStarted droplet and container started
	ModeStarted

	// ModeReady server is running and can be joined
	ModeReady

	// ModeShuttingDown server is currently being shut down
	ModeShuttingDown
)
