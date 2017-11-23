package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
)

var bastion *gobastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	utils.Send(w, res)
}

func main() {
	bastion = gobastion.NewBastion()
	bastion.APIRouter.Get("/hello", helloHandler)
	bastion.Run(":8080")
}
