package main

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
)

func main() {
	app := bastion.New(nil)
	reader := new(bastion.JsonReader)
	handler := todo.Handler{
		Reader:    reader,
		Responder: bastion.DefaultResponder,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
