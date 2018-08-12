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
	app.APIRouter.Get("/recovery", handler)
	app.Serve()
}
