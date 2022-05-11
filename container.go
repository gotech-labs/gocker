package gocker

import (
	dockertest "github.com/ory/dockertest/v3"
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
	return c.resource.GetHostPort(id)
}

func (c *Container) Purge() {
	Purge(c.resource)
}
