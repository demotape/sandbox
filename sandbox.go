package sandbox

import (
	"../banner"
	"../config"
	"../ssh-key"
	"../tutorial"
	"./bind_volume"
	"./port_mapping"
	"./runtime_env"
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

// Taken from/proc/sys/net/ipv4/ip_local_port_range
const (
	PortMin        = 32768
	PortMax        = 61000
	secretKey      = "/tmp/demotape_ssh_secret"
	sessionDirBase = "/tmp/demotape"
)

type Sandbox struct {
	ImageName   string                  `json:"image_name"`
	RuntimeEnv  *runtime_env.RuntimeEnv `json:"runtime_env"`
	ContainerId string                  `json:"container_id"`
	BindVolume  bind_volume.BindVolume
	PortMapping port_mapping.PortMapping
}

var (
	client *docker.Client
	HostIP = "192.168.33.10"
)

func init() {
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		panic("DOCKER_HOST is not set")
	}
	client, _ = docker.NewClient(endpoint)

	os.MkdirAll(sessionDirBase, 644)
}

func (s Sandbox) Start() (string, string, error) {
	imageName := s.ImageName
	bindVol := s.BindVolume

	var err error
	var container *docker.Container

	containerConfig := docker.Config{
		AttachStdout: true,
		AttachStdin:  true,
		Image:        imageName,
		ExposedPorts: s.PortMapping.ExposedPorts,
	}

	opt := docker.CreateContainerOptions{Config: &containerConfig}

	container, err = client.CreateContainer(opt)
	if err != nil {
		fmt.Println("Fail to create container")
		return "", "", err
	}

	hostConfig := docker.HostConfig{
		PortBindings: s.PortMapping.PortBindings,
		Binds:        bindVol.VolumeMappings(),
	}

	err = client.StartContainer(container.ID, &hostConfig)

	if err != nil {
		fmt.Printf("Fail to start container %s", container.ID)
		return "", "", err
	} else {
		fmt.Println(container.ID)
	}

	inspect, _ := client.InspectContainer(container.ID)

	return inspect.NetworkSettings.IPAddress, container.ID, nil
}

func (s *Sandbox) Stop() error {
	return client.StopContainer(s.ContainerId, 0)
}

func (s *Sandbox) Commit() error {
	// containerOpt := docker.CommitContainerOptions{Container: s.ContainerId, Repository: name}
	containerOpt := docker.CommitContainerOptions{Container: s.ContainerId}
	_, err := client.CommitContainer(containerOpt)

	if err != nil {
		return err
	}

	return nil
}

func (s *Sandbox) LoginCommands() []string {
	r, _ := regexp.Compile("\n$")
	keyContent := string(s.RuntimeEnv.SshEnv.KeyPair.Secret[:len(s.RuntimeEnv.SshEnv.KeyPair.Secret)])
	keyContent = r.ReplaceAllString(keyContent, "")

	return []string{
		"cat << EOF > " + secretKey,
		keyContent,
		"EOF",
		"chmod 600 " + secretKey,
		"ssh -o StrictHostKeychecking=no -o UserKnownHostsFile=/dev/null -p " + strconv.Itoa(s.RuntimeEnv.SshEnv.PortNum) + " " + "root@" + HostIP + " -i " + secretKey,
	}
}

func BuildImage(imageName string, tag string, dockerfile string) error {
	taggedImage := imageName + ":" + tag
	buildOpt := docker.BuildImageOptions{
		Name:           taggedImage,
		SuppressOutput: true,
		RmTmpContainer: true,
		InputStream:    CreateDockerfile(dockerfile),
		OutputStream:   bytes.NewBuffer(nil),
	}

	err := client.BuildImage(buildOpt)
	if err != nil {
		fmt.Println("Failed to build base image")
		panic(err)
	}
	return nil
}

func CreateDockerfile(content string) io.Reader {
	size := int64(len(content))
	t := time.Now()
	inputbuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: size, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte(content))
	tr.Close()
	return inputbuf
}

func CreateRunTimeEnv(tutorial *tutorial.Tutorial, name string) (*runtime_env.RuntimeEnv, error) {
	runtimeDir, err := runtime_env.CreateRuntimeDir()

	if err != nil {
		fmt.Println("Fail to create runtime dir", err)
		return nil, err
	}

	sshEnv, err := prepareSsh(runtimeDir)

	if err != nil {
		fmt.Println("Fail to prepareSsh")
		return nil, err
	}

	prepareWelcome(runtimeDir, tutorial, name)

	runtimeEnv := runtime_env.RuntimeEnv{RuntimeDir: runtimeDir, SshEnv: sshEnv}

	return &runtimeEnv, nil
}

func prepareWelcome(runtimeDir string, tutorial *tutorial.Tutorial, name string) {
	welcomeDir := path.Join(runtimeDir, config.WelcomeDir)
	os.MkdirAll(welcomeDir, 0644)

	tutoFile := path.Join(welcomeDir, "tutorial.sh")
	ioutil.WriteFile(tutoFile, []byte(tutorial.ToCommands()), 0766)

	bannerFile := path.Join(welcomeDir, "banner.sh")
	ioutil.WriteFile(bannerFile, []byte(banner.Welcome(name)), 0766)
}

func prepareSsh(runtimeDir string) (*runtime_env.SshEnv, error) {
	sessionDir, err := createSessionDir("build")

	keyPair, err := sshkey.GenerateKeys(sessionDir)
	if err != nil {
		return nil, err
	}

	dotSshDir := path.Join(runtimeDir, ".ssh")
	os.Mkdir(dotSshDir, 0644)

	authorizedKeys := path.Join(dotSshDir, "authorized_keys")
	ioutil.WriteFile(authorizedKeys, keyPair.Pub, 0400)

	sshEnv := runtime_env.SshEnv{randomPortNum(), keyPair, dotSshDir}

	return &sshEnv, nil
}

func randomPortNum() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(PortMax-PortMin) + PortMin
}

func createSessionDir(prefix string) (string, error) {
	dir, err := ioutil.TempDir(sessionDirBase, prefix)

	if err != nil {
		fmt.Println("Fail to create tempdir")
		return "", err
	}

	return dir, nil
}

type ErrorSandboxNotExists struct {
	Msg    string
	Offset int64
}

func (e ErrorSandboxNotExists) Error() string { return e.Msg }
