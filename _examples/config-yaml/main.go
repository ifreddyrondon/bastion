package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/config"
	"github.com/ifreddyrondon/bastion/render/json"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	json.NewRenderer(w).Send(res)
}

func main() {
	cfg, _ := config.FromFile("./config.yaml")
	app := bastion.New(cfg)
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
