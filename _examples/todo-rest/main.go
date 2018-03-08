package main

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
	jsonreader "github.com/ifreddyrondon/bastion/reader/json"
	jsonrender "github.com/ifreddyrondon/bastion/renderer/json"
)

func main() {
	app := bastion.New(nil)
	handler := todo.Handler{
		Reader: jsonreader.NewReader,
		Render: jsonrender.NewRenderer,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	app.Serve()
}
