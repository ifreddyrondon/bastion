package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/config"
)

var app *gobastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

func main() {
	cfg, _ := config.FromFile("./config.yaml")
	app = gobastion.New(cfg)
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
