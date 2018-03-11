package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	json.NewRender(w).Send(res)
}

func main() {
	app, _ := bastion.FromFile("./options.yaml")
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
