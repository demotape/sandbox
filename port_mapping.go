package sandbox

import (
	"github.com/fsouza/go-dockerclient"
)

type PortMapping struct {
	ExposedPorts map[docker.Port]struct{}
	PortBindings map[docker.Port][]docker.PortBinding
}

//Ex. NewPortMapping("22/tcp", "42456")
func (pm *PortMapping) AddBinding(exposedPorts docker.Port, bindPorts string) {
	//pm := PortMapping{}
	pm.ExposedPorts = make(map[docker.Port]struct{})
	pm.ExposedPorts[docker.Port(exposedPorts)] = struct{}{}

	pm.PortBindings = make(map[docker.Port][]docker.PortBinding)
	pm.PortBindings[exposedPorts] = []docker.PortBinding{{
		HostIp:   "0.0.0.0",
		HostPort: bindPorts,
	}}

}
