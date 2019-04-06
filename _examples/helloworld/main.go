package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, r *http.Request) {
	l := bastion.LoggerFromCtx(r.Context())
	l.Info().Msg("handler")

	res := struct {
		Message string `json:"message"`
	}{"hello world"}
	render.NewJSON().Send(w, res)
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	fmt.Fprintln(os.Stderr, app.Serve())
}
