package types

// State stores information about the CS:GO server
type State struct {
	DropletID   int    `json:"dropletId"`
	DropletIP   string `json:"dropletIp"`
	ContainerID string `json:"containerId"`
	Progress    int    `json:"progress"`
	Mode        `json:"mode"`
}
