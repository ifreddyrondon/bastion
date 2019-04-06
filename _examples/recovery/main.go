package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	panic("Something wrong happened")
}

func main() {
	app := bastion.New()
	app.Get("/recovery", handler)
	fmt.Fprintln(os.Stderr, app.Serve())
}
