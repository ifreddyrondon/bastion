package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	render.JSON.Send(w, map[string]string{"message": "hello bastion"})
}

func onShutdown() {
	log.Printf("My registered on shutdown. Doing something...")
}

func main() {
	app := bastion.New()
	app.RegisterOnShutdown(onShutdown)
	app.Get("/hello", helloHandler)
	fmt.Fprintln(os.Stderr, app.Serve())
}
