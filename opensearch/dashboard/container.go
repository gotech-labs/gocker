package dashboard

import (
	"errors"
	"io"
	"net/http"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"

	"github.com/gotech-labs/gocker"
)

func New(tag, opensearchHost string) *Container {
	dockerOptions := []gocker.ConfigOption{
		gocker.WithEnv(gocker.Env{
			"OPENSEARCH_HOSTS":                   opensearchHost,
			"DISABLE_SECURITY_DASHBOARDS_PLUGIN": true,
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
					if string(body) == "OpenSearch Dashboards server is not ready yet" {
						err = errors.New(string(body))
					}
				}
			}
			return err
		}),
	}
	return &Container{
		Container: gocker.New(
			"opensearch-dashboards_"+tag,
			"opensearchproject/opensearch-dashboards",
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
