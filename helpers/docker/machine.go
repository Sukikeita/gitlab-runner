package docker_helpers

import (
	"time"
)

type Machine interface {
	Create(driver, name string, opts ...string) error
	Provision(name string) error
	Remove(name string) error
	Stop(name string, timeout time.Duration) error
	List() (machines []string, err error)
	Exist(name string) bool

	CanConnect(name string) bool
	Credentials(name string) (DockerCredentials, error)
}
