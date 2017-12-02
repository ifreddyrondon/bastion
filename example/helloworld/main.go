package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
)

var app *gobastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	utils.Send(w, res)
}

func main() {
	app = gobastion.New("")
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
