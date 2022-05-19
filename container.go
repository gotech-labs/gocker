package gocker

import (
	"net/url"
	"os"

	dockertest "github.com/ory/dockertest/v3"

	"github.com/gotech-labs/core/log"
)

type Container struct {
	name     string
	tag      string
	resource *dockertest.Resource
}

func (c *Container) Tag() string {
	return c.tag
}

func (c *Container) Name() string {
	return c.name
}

func (c *Container) IPAddress() string {
	return c.resource.Container.NetworkSettings.IPAddress
}

func (c *Container) Port(id string) string {
	return c.resource.GetPort(id)
}

func (c *Container) HostPort(id string) string {
	if dockerURL := os.Getenv("DOCKER_HOST"); dockerURL != "" {
		u, err := url.Parse(dockerURL)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to parse docker url")
		}
		return u.Hostname() + ":" + c.Port(id)
	}
	return c.resource.GetHostPort(id)
}

func (c *Container) Purge() {
	Purge(c.resource)
}
