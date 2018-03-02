package main

import (
	"log"
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

var app *bastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

func onShutdown() {
	log.Printf("My registered on shutdown. Doing something...")
}

func main() {
	app = bastion.New(nil)
	app.RegisterOnShutdown(onShutdown)
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
