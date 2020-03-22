package server

// State stores information about the CS:GO server
type State struct {
	DropletID   int    `json:"dropletId"`
	DropletIP   string `json:"dropletIp"`
	ContainerID string `json:"containerId"`
	Mode        `json:"mode"`
}
