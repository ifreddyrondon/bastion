package main

import (
	"fmt"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
)

func main() {
	app := bastion.New()
	app.Mount("/todo/", todo.Routes())
	fmt.Fprintln(os.Stderr, app.Serve())
}
