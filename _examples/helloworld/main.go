package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	render.JSONRender(w).Send(res)
}

func main() {
	app := bastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
