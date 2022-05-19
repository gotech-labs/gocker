package gocker

import (
	"fmt"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"
)

type config struct {
	runOptions     *dockertest.RunOptions
	awaitRetryFunc func(*dockertest.Resource) error
	hostConfigFunc func(*docker.HostConfig)
}

type ConfigOption func(*config)
type Env map[string]interface{}

func WithEnv(env Env) ConfigOption {
	return func(cfg *config) {
		values := make([]string, 0)
		for k, v := range env {
			values = append(values, fmt.Sprintf("%s=%s", k, fmt.Sprint(v)))
		}
		cfg.runOptions.Env = values
	}
}

func WithCmd(cmd ...string) ConfigOption {
	return func(cfg *config) {
		cfg.runOptions.Cmd = cmd
	}
}

func WithEntryPoint(entryPoint ...string) ConfigOption {
	return func(cfg *config) {
		cfg.runOptions.Entrypoint = entryPoint
	}
}

func WithPortBindings(bindigPorts ...string) ConfigOption {
	return func(cfg *config) {
		for _, port := range bindigPorts {
			var (
				key   = docker.Port(port + "/tcp")
				value = []docker.PortBinding{{HostPort: port}}
			)
			if cfg.runOptions.PortBindings == nil {
				cfg.runOptions.PortBindings = map[docker.Port][]docker.PortBinding{key: value}
			} else {
				cfg.runOptions.PortBindings[key] = value
			}
		}
	}
}

func WithAwaitRetryFunc(retryFunc func(*dockertest.Resource) error) ConfigOption {
	return func(cfg *config) {
		cfg.awaitRetryFunc = retryFunc
	}
}

func WithHostConfigFunc(configFunc func(*docker.HostConfig)) ConfigOption {
	return func(cfg *config) {
		cfg.hostConfigFunc = configFunc
	}
}
