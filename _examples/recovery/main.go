package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	panic("Something wrong happened")
}

func main() {
	app := bastion.New()
	app.Get("/recovery", handler)
	app.Serve()
}
