package sandbox

import (
	"io/ioutil"
	"os"
)

type RuntimeEnv struct {
	RuntimeDir string
	SshEnv     *SshEnv
}

type SshEnv struct {
	PortNum   int
	KeyPair   *KeyPair
	DotSshDir string
}

func CreateRuntimeDir() (string, error) {
	os.MkdirAll(Host.RunTimeDir, 0644)
	dir, err := ioutil.TempDir(Host.RunTimeDir, "runtime")

	if err != nil {
		return "", err
	}

	return dir, nil
}
