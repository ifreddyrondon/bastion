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
	render.JSON.Send(w, map[string]string{"message": "hello bastion"})
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	fmt.Fprintln(os.Stderr, app.Serve())
}
