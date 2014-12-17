package demotape

import (
	"path"
)

var (
	ContainerTmpDir = "/var/tmp"
	WelcomeDir      = "/welcome"
)

type ContainerConfig struct {
	TmpDir     string
	WelcomeDir string
}

type HostConfig struct {
	RunTimeDir string
}

var Container = ContainerConfig{
	TmpDir:     ContainerTmpDir,
	WelcomeDir: path.Join(ContainerTmpDir, WelcomeDir),
}

var Host = HostConfig{
	RunTimeDir: "/var/lib/demotape",
}
