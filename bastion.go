package gobastion

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an api router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	Router *chi.Mux
}

// NewRouter returns a new GogApp instance ready
func NewBastion() *Bastion {
	app := new(Bastion)
	app.Router = chi.NewRouter()
	initialize(app.Router)
	return app
}

func initialize(r *chi.Mux) {
	r.Get("/ping", pingHandler)
}

func (app *Bastion) Run(address string) {
	log.Printf("Running on %s", address)
	log.Fatal(http.ListenAndServe(address, app.Router))
}
