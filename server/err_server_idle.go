package server

// ErrServerIdle states the server is idle
type ErrServerIdle struct{}

func (err ErrServerIdle) Error() string {
	return "server is idle"
}
