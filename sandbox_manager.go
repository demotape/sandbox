package main

import (
	"../../go-demotape/controller"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/codegangsta/negroni"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Demotape sandbox manager started")
	mux := pat.New()

	mux.Post("/sandboxes", http.HandlerFunc(controller.CreateSandbox))
	mux.Put("/sandboxes/:image_name/start", http.HandlerFunc(controller.StartSandbox))
	mux.Put("/sandboxes/:container_id/checkin", http.HandlerFunc(controller.CheckinSandbox))

	n := negroni.Classic()
	n.UseHandler(mux)

	if os.Getenv("DEMOTAPE_ENV") == "test" {
		n.Run(":3003")
	} else {
		n.Run(":3000")
	}
}
