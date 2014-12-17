package demotape

import (
	"../../config"
	"../../ssh-key"
	"io/ioutil"
	"os"
)

type RuntimeEnv struct {
	RuntimeDir string
	SshEnv     *SshEnv
}

type SshEnv struct {
	PortNum   int
	KeyPair   *sshkey.KeyPair
	DotSshDir string
}

func CreateRuntimeDir() (string, error) {
	os.MkdirAll(config.Host.RunTimeDir, 0644)
	dir, err := ioutil.TempDir(config.Host.RunTimeDir, "runtime")

	if err != nil {
		return "", err
	}

	return dir, nil
}
