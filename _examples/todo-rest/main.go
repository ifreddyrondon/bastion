package main

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
	"github.com/ifreddyrondon/bastion/render/json"
)

func main() {
	app := bastion.New(nil)
	reader := new(bastion.JsonReader)
	handler := todo.Handler{
		Reader: reader,
		Render: json.NewRenderer,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
