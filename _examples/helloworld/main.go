package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func handler(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	l := bastion.LoggerFromCtx(r.Context())
	l.Info().Msg("handler")

	json.NewRender(w).Send(res)
}

func main() {
	app := bastion.New(bastion.Options{})
	app.APIRouter.Get("/hello", handler)
	app.Logger.Info().Str("app", "test").Msg("main")
	app.Serve()
}
