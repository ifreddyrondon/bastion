package gobastion

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	r         *chi.Mux
	APIRouter chi.Router
}

// NewRouter returns a new GogApp instance ready
func NewBastion() *Bastion {
	app := new(Bastion)
	initialize(app)
	return app
}

func (app *Bastion) Run(address string) {
	log.Printf("Running on %s", address)
	log.Fatal(http.ListenAndServe(address, app.r))
}

func initialize(app *Bastion) {
	/**
	 * internal router
	 */
	app.r = chi.NewRouter()

	/**
	 * Ping route
	 */
	app.r.Get("/ping", pingHandler)

	/**
	 * API Router
	 */
	app.APIRouter = chi.NewRouter()
	app.APIRouter.Use(middleware.RequestID)
	app.APIRouter.Use(middleware.Logger)
	app.r.Mount("/", app.APIRouter)
}
