package gobastion

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ifreddyrondon/gobastion/config"
)

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	r         *chi.Mux
	cfg       *config.Config
	APIRouter chi.Router
}

// NewRouter returns a new Bastion instance.
// if configPath is empty the configuration will be from defaults.
// 	Defaults:
//		api:
//			base_path: "/"
// Otherwise the configuration will be loaded from configPath.
// If the config file is missing or unable to unmarshal the will panic.
func NewBastion(configPath string) *Bastion {
	app := new(Bastion)
	app.cfg = config.New()
	if configPath != "" {
		if err := app.cfg.FromFile(configPath); err != nil {
			log.Panic(err)
		}
	}
	initialize(app)
	return app
}

// Serve the application at the specified address/port
func (app *Bastion) Serve(address string) {
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
	app.r.Mount(app.cfg.API.BasePath, app.APIRouter)
}
