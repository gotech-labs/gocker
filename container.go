package gocker

import (
	dockertest "github.com/ory/dockertest/v3"
)

type Container struct {
	resource *dockertest.Resource
}

func (c *Container) Port(id string) string {
	return c.resource.GetPort(id)
}

func (c *Container) Purge() {
	Purge(c.resource)
}
