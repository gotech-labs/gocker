package minio

import (
	"fmt"
	"net/http"

	dockertest "github.com/ory/dockertest/v3"

	"github.com/gotech-labs/gocker"
)

func New(tag string, bucket string) *Container {
	dockerOptions := []gocker.ConfigOption{
		gocker.WithEnv(gocker.Env{
			"MINIO_ROOT_USER":     credentials.User,
			"MINIO_ROOT_PASSWORD": credentials.Password,
			//"MINIO_BROWSER_REDIRECT_URL": "http://localhost:9001",
		}),
		gocker.WithCmd(
			"server", "/data", "--console-address", ":9001",
		),
		gocker.WithAwaitRetryFunc(func(resource *dockertest.Resource) error {
			address := "http://" + resource.GetHostPort("9000/tcp") + "/minio/health/live"
			println(address)
			resp, err := http.Get(address)
			if err == nil && resp.ContentLength > 0 {
				err = resp.Body.Close()
			}
			return err
		}),
	}
	return &Container{
		Container: gocker.New(
			"minio_"+tag,
			"minio/minio",
			tag,
			dockerOptions...,
		),
	}
}

type Container struct {
	*gocker.Container
}

func (c *Container) RestAPIEndpoint() string {
	return fmt.Sprintf("http://%s", c.HostPort("9000/tcp"))
}

func (c *Container) ConsoleAddress() string {
	return fmt.Sprintf("http://%s", c.HostPort("9001/tcp"))
}

var (
	credentials = struct {
		User     string
		Password string
	}{
		User:     "minio_user",
		Password: "minio_password",
	}
)
