package server

// ErrServerStarted states the server already started
type ErrServerStarted struct{}

func (err ErrServerStarted) Error() string {
	return "server already started"
}
