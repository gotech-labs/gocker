package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	dockertest "github.com/ory/dockertest/v3"
	docker "github.com/ory/dockertest/v3/docker"

	"github.com/gotech-labs/gocker"
)

func New(tag string, dbName string) *Container {
	dockerOptions := []gocker.ConfigOption{
		gocker.WithEnv(gocker.Env{
			"MYSQL_ROOT_PASSWORD":        rootPassword,
			"MYSQL_ALLOW_EMPTY_PASSWORD": "N",
			"MYSQL_DATABASE":             dbName,
			"MYSQL_USER":                 user,
			"MYSQL_PASSWORD":             password,
		}),
		gocker.WithHostConfigFunc(func(hostConfig *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			hostConfig.AutoRemove = true
			hostConfig.RestartPolicy = docker.RestartPolicy{Name: "no"}
		}),
		gocker.WithAwaitRetryFunc(func(resource *dockertest.Resource) error {
			url := databaseURL(resource.GetHostPort("3306/tcp"), dbName)
			db, err := sql.Open("mysql", url)
			if err == nil {
				defer db.Close()
				err = db.Ping()
			}
			return err
		}),
	}
	return &Container{
		Container: gocker.New(
			"mysql_"+tag,
			"mysql",
			tag,
			dockerOptions...,
		),
		dbName: dbName,
	}
}

type Container struct {
	*gocker.Container
	dbName string
}

func (c *Container) DatabaseURL() string {
	return databaseURL(c.HostPort("3306/tcp"), c.dbName)
}

var (
	user         = "docker"
	password     = "docker"
	rootPassword = "password"
)

func databaseURL(hostPort, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		hostPort,
		dbName,
	)
}
