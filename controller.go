package sandbox

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

type OkResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SandboxStarted struct {
	LoginCommands []string `json:"login_commands"`
	ContainerId   string   `json:"container_id"`
}

type StartSandboxData struct {
	ImageName string `json:"image_name"`
	DemoName  string `json:"demo_name"`
	Author    string `json:"author"`
	Tutorial  Tutorial
}

type CreateSandboxData struct {
	ImageName string `json:"image_name"`
}

type SandboxCreated struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Error string `json:"error"`
}

func Respond(w http.ResponseWriter, httpStatus int, msg interface{}) {
	rd := jsonRender()
	rd.JSON(w, httpStatus, msg)
}
func jsonRender() *render.Render {
	return render.New(render.Options{
		IndentJSON: true,
	})
}

func CreateSandbox(w http.ResponseWriter, req *http.Request) {
	var d CreateSandboxData
	var res interface{}
	var err interface{}

	json.NewDecoder(req.Body).Decode(&d)
	err = createSandbox(d.ImageName)

	if err != nil {
		log.Printf("Got unknown error type: %s\n", err)
		res = Error{
			Error: "Unknown error",
		}
		Respond(w, 500, res)
		return
	} else {
		Respond(w, 201, nil)

		return
	}
}

func StartSandbox(w http.ResponseWriter, r *http.Request) {
	var data StartSandboxData
	var res interface{}

	params := r.URL.Query()
	name := params.Get(":image_name")
	json.NewDecoder(r.Body).Decode(&data)

	sandbox, err := startSandbox(data.ImageName, data.DemoName, data.Author, &data.Tutorial)

	if sandbox != nil && err == nil {
		res = SandboxStarted{
			LoginCommands: sandbox.LoginCommands(),
			ContainerId:   sandbox.ContainerId,
		}
		Respond(w, 200, res)
		return
	} else if sandbox == nil && err == nil {
		res = Error{
			Error: fmt.Sprintf("%s not found", name),
		}
		Respond(w, 404, res)
		return
	} else {
		fmt.Println("demo", sandbox, "err", err)
		res = Error{
			Error: "Server Error",
		}
		Respond(w, 505, res)
		return
	}
}

func CheckinSandbox(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query()
	name := param.Get(":container_id")
	var res interface{}

	err := checkinSandbox(name)

	if err != nil {
		res = Error{
			Error: fmt.Sprintf("Fail to checkin.\nError: %s", err.Error()),
		}
		Respond(w, 505, res)
		return
	}

	Respond(w, 200, nil)
	return
}

func createSandbox(name string) error {

	dockerfile := fmt.Sprintf("From %s\n", "demotape")
	tag := "latest"
	err := BuildImage(name, tag, dockerfile)

	if err != nil {
		fmt.Println("Fail to create base image")
		return err
	}

	return nil
}

func startSandbox(name string, demoName string, author string, tutorial *Tutorial) (*Sandbox, error) {

	runtimeEnv, err := CreateRunTimeEnv(tutorial, demoName, author)
	bindVol := NewBindVolumes(runtimeEnv)

	if err != nil {
		fmt.Println("Fail to prepareSsh", err)
		return nil, err
	}

	sshEnv := runtimeEnv.SshEnv

	portMapping := &PortMapping{}
	portMapping.AddBinding("22/tcp", strconv.Itoa(sshEnv.PortNum))

	s := Sandbox{
		ImageName:   fmt.Sprintf("%s:latest", name),
		RuntimeEnv:  runtimeEnv,
		BindVolume:  bindVol,
		PortMapping: *portMapping,
	}

	_, id, err := s.Start()

	if err != nil {
		fmt.Println("Fail to run container")
		return nil, err
	}

	s.ContainerId = id

	return &s, nil
}

func checkinSandbox(containerId string) error {
	var err error
	s := Sandbox{ContainerId: containerId}

	err = s.Commit()

	if err != nil {
		return err
	}

	err = s.Stop()

	if err != nil {
		return err
	}

	return nil
}

// TODO: move this to util package
func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	if pathError, ok := err.(*os.PathError); ok {
		if pathError.Err == syscall.ENOTDIR {
			return false
		}
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// TODO: move this to util package
func randomSha() string {
	s := strconv.FormatInt(time.Now().Unix(), 10)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
