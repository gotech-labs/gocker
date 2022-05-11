package opensearch

import (
	"fmt"
	"net/http"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"

	"github.com/gotech-labs/gocker"
)

func New(tag string, plugins ...string) *Container {
	dockerOptions := []gocker.ConfigOption{
		gocker.WithEnv(gocker.Env{
			"cluster.name":                "docker-cluster",
			"node.name":                   "docker-node",
			"discovery.type":              "single-node",
			"bootstrap.memory_lock":       true,
			"network.host":                "0.0.0.0",
			"OPENSEARCH_JAVA_OPTS":        "-Xms1024m -Xmx1024m",
			"DISABLE_INSTALL_DEMO_CONFIG": true,
			"DISABLE_SECURITY_PLUGIN":     true,
		}),
		gocker.WithHostConfigFunc(func(hostConfig *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			hostConfig.AutoRemove = true
			hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
			hostConfig.Ulimits = []docker.ULimit{
				{Name: "memlock", Soft: -1, Hard: -1},
				{Name: "nofile", Soft: 65535, Hard: 65535},
			}
		}),
		gocker.WithAwaitRetryFunc(func(resource *dockertest.Resource) error {
			endpoint := "http://" + resource.GetHostPort("9200/tcp")
			resp, err := http.Get(endpoint)
			if err == nil && resp.ContentLength > 0 {
				err = resp.Body.Close()
			}
			return err
		}),
	}
	if len(plugins) > 0 {
		cmd := ""
		for _, plugin := range plugins {
			cmd += "${OPENSEARCH_HOME}/bin/opensearch-plugin install " + plugin + ";"
		}
		dockerOptions = append(dockerOptions, gocker.WithCmd(
			"/bin/sh",
			"-c",
			cmd+"./opensearch-docker-entrypoint.sh"))
	}
	return &Container{
		Container: gocker.New(
			"opensearch_"+tag,
			"opensearchproject/opensearch",
			tag,
			dockerOptions...,
		),
	}
}

type Container struct {
	*gocker.Container
}

func (c *Container) RestAPIEndpoint() string {
	return fmt.Sprintf("http://%s", c.HostPort("9200/tcp"))
}

func (c *Container) AcceptDashboardURL() string {
	return fmt.Sprintf("http://%s:%s", c.Name(), "9200")
}
