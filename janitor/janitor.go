package janitor

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
	"time"
)

var (
	client   *docker.Client
	endpoint string
)

func init() {
	endpoint = os.Getenv("DOCKER_HOST")

	fmt.Println(endpoint)

	if endpoint == "" {
		panic("DOCKER_HOST is not set")
	}
}

func Clean() int {
	timeLimit := int64(1000)
	client, _ := docker.NewClient(endpoint)
	currentTime := time.Now().Local().Unix()

	option := docker.ListContainersOptions{}

	containers, _ := client.ListContainers(option)

	for _, container := range containers {
		createdDate := container.Created
		runningTime := currentTime - createdDate

		if runningTime > timeLimit {
			client.StopContainer(container.ID, 0)
			fmt.Printf("Kill %s\n", container.ID)
		}
	}

	fmt.Println("janitor is called!!!")

	return 1

}
