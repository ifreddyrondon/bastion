package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
)

var app *gobastion.Bastion

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

func main() {
	app = gobastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
