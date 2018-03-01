package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

var app *bastion.Bastion

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

func main() {
	app = bastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
