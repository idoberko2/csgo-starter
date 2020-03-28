package types

// DockerCreator represents entities that generate Docker API clients
type DockerCreator interface {
	Create(ip string) (Docker, error)
}
