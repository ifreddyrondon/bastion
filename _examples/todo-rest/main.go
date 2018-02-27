package main

import (
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/_examples/todo-rest/todo"
)

func main() {
	app := gobastion.New(nil)
	reader := new(gobastion.JsonReader)
	handler := todo.Handler{
		Reader:    reader,
		Responder: gobastion.DefaultResponder,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
