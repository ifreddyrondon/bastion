package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	render.NewJSON().Send(w, res)
}

func main() {
	app, _ := bastion.FromFile("./options.yaml")
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
