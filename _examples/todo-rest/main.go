package main

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
	"github.com/ifreddyrondon/bastion/render"
)

func main() {
	app := bastion.New()
	handler := todo.Handler{
		Render: render.NewJSON(),
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
