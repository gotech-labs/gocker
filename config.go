package gocker

import (
	"fmt"

	dockertest "github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

type config struct {
	runOptions     *dockertest.RunOptions
	awaitRetryFunc func(*dockertest.Resource) error
	hostConfigFunc func(*dc.HostConfig)
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

func WithAwaitRetryFunc(retryFunc func(*dockertest.Resource) error) ConfigOption {
	return func(cfg *config) {
		cfg.awaitRetryFunc = retryFunc
	}
}

func WithHostConfigFunc(configFunc func(*dc.HostConfig)) ConfigOption {
	return func(cfg *config) {
		cfg.hostConfigFunc = configFunc
	}
}
