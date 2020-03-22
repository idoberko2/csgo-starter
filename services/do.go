package services

import (
	"context"
	"errors"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/digitalocean/godo"

	log "github.com/sirupsen/logrus"
)

// Do handles interactions with DigitalOcean
type Do struct {
	client *godo.Client
}

// NewDo creates a new instance of Do
func NewDo(token string) *Do {
	return &Do{
		client: godo.NewFromToken(token),
	}
}

// StartDroplet starts a new Droplet adjusted for CS:GO
func (do *Do) StartDroplet(ctx context.Context) (string, int, error) {
	var (
		ip  string
		err error
	)
	createRequest := &godo.DropletCreateRequest{
		Name:   "cs-go-droplet",
		Region: "fra1",
		Size:   "s-2vcpu-4gb",
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{
				Fingerprint: os.Getenv("DO_SSH_KEY_FP"),
			},
		},
		Image: godo.DropletCreateImage{
			Slug: "docker-18-04",
		},
		// required initialization of the droplet
		UserData: `#cloud-config

runcmd:
    - docker run -d -p 2375:2375 -v /var/run/docker.sock:/var/run/docker.sock jarkt/docker-remote-api
    - sudo ufw allow 27015/tcp
    - sudo ufw allow 27015/udp
`,
	}

	newDroplet, _, err := do.client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return "", 0, err
	}
	dropletID := newDroplet.ID

	ip, err = do.waitForIP(ctx, dropletID)
	if err != nil {
		return "", 0, err
	}

	return ip, dropletID, nil
}

func (do *Do) waitForIP(ctx context.Context, dropletID int) (string, error) {
	var ip string

	err := backoff.Retry(func() error {
		var err error
		droplet, _, err := do.client.Droplets.Get(ctx, dropletID)
		log.Debug("Checking IP")
		ip, err = droplet.PublicIPv4()
		if err != nil {
			return err
		}
		if ip == "" {
			log.Debug("IP is empty")
			return errors.New("ip is empty")
		}
		log.WithField("ip", ip).Debug("Received IP")

		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return "", err
	}

	return ip, nil
}

// StopDroplet stops the droplet identified by dropletID
func (do *Do) StopDroplet(ctx context.Context, dropletID int) error {
	_, err := do.client.Droplets.Delete(ctx, dropletID)
	if err != nil {
		return err
	}

	return nil
}
