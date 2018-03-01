package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/config"
)

var app *bastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

func main() {
	cfg, _ := config.FromFile("./config.yaml")
	app = bastion.New(cfg)
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
