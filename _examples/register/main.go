package main

import (
	"log"
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

func onShutdown() {
	log.Printf("My registered on shutdown. Doing something...")
}

func main() {
	app := bastion.New(bastion.Options{})
	app.RegisterOnShutdown(onShutdown)
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
