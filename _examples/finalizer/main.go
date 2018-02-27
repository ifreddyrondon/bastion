package main

import (
	"net/http"

	"log"

	"github.com/ifreddyrondon/gobastion"
)

var app *gobastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	app.Send(w, res)
}

type MyFinalizer struct{}

func (f MyFinalizer) Finalize() error {
	log.Printf("[finalizer:MyFinalizer] doing something")
	return nil
}

func main() {
	app = gobastion.New(nil)
	app.AppendFinalizers(MyFinalizer{})
	app.APIRouter.Get("/hello", helloHandler)
	app.Serve()
}
