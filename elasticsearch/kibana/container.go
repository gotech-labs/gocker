package kibana

import (
	"errors"
	"io"
	"net/http"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"

	"github.com/gotech-labs/gocker"
)

func New(tag, esHost string) *Container {
	dockerOptions := []gocker.ConfigOption{
		gocker.WithEnv(gocker.Env{
			"ELASTICSEARCH_HOSTS": esHost,
		}),
		gocker.WithHostConfigFunc(func(hostConfig *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			hostConfig.AutoRemove = true
		}),
		gocker.WithAwaitRetryFunc(func(resource *dockertest.Resource) error {
			endpoint := "http://" + resource.GetHostPort("5601/tcp")
			resp, err := http.Get(endpoint)
			if err == nil && resp.ContentLength > 0 {
				var body []byte
				if body, err = io.ReadAll(resp.Body); err == nil {
					if string(body) == "Kibana server is not ready yet" {
						err = errors.New(string(body))
					}
				}
			}
			return err
		}),
	}
	return &Container{
		Container: gocker.New(
			"kibana_"+tag,
			"kibana",
			tag,
			dockerOptions...,
		),
	}
}

type Container struct {
	*gocker.Container
}

func (c *Container) RestAPIEndpoint() string {
	return "http://" + c.HostPort("5601/tcp")
}
