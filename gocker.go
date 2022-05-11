package gocker

import (
	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"

	"github.com/gotech-labs/core/log"
)

func New(name, repository, tag string, options ...ConfigOption) *Container {
	var (
		resource *dockertest.Resource
		exist    bool
		err      error
	)
	if resource, exist = pool.ContainerByName(name); !exist {
		log.Info().Msgf("Startup docker resource [%s]", name)
		cfg := &config{
			runOptions: &dockertest.RunOptions{
				Name:       name,
				Repository: repository,
				Tag:        tag,
			},
			awaitRetryFunc: func(*dockertest.Resource) error { return nil },
			hostConfigFunc: func(*docker.HostConfig) {},
		}
		for _, option := range options {
			option(cfg)
		}
		// pulls an image, creates a container based on it and runs it
		resource, err = pool.RunWithOptions(cfg.runOptions, cfg.hostConfigFunc)
		if err != nil {
			log.Panic().Err(err).Msgf("Could not start resource %s: %s", name, err)
		}
		// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
		retry := 0
		if err := pool.Retry(func() error {
			err := cfg.awaitRetryFunc(resource)
			if err != nil {
				retry++
				log.Info().Msgf("Waiting for startup %s [retry=%d, err=%s]", name, retry, err)
			} else {
				log.Info().Msgf("Ready to accept resource %s", name)
			}
			return err
		}); err != nil {
			Purge(resource)
			log.Panic().Err(err).Msgf("Could not connect to resource %s: %s", name, err)
		}
	}
	return &Container{
		resource: resource,
	}
}

func Purge(resource *dockertest.Resource) {
	if resource != nil {
		if err := pool.Purge(resource); err != nil {
			log.Panic().Err(err).Msgf("Could not purge resource %s: %s", resource.Container.Name, err)
		}
	}
}

var pool = func() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panic().Err(err).Msgf("Could not connect to docker: %s", err)
	}
	return pool
}()
