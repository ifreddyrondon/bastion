package main

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
	"github.com/ifreddyrondon/bastion/renderer/json"
)

func main() {
	app := bastion.New(nil)
	handler := todo.Handler{
		Render: json.NewRenderer,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
