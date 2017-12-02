package gobastion

import (
	"context"
	"log"
	"net/http"

	"os"

	"syscall"

	"github.com/go-chi/chi"
	CHIMiddleware "github.com/go-chi/chi/middleware"
	"github.com/ifreddyrondon/gobastion/config"
	"github.com/ifreddyrondon/gobastion/midleware"
	"github.com/markbates/sigtx"
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

// New returns a new Bastion instance.
// if configPath is empty the configuration will be from defaults.
// 	Defaults:
//		api:
//			base_path: "/"
//		server:
//			address ":8080"
// Otherwise the configuration will be loaded from configPath.
// If the config file is missing or unable to unmarshal the will panic.
func New(configPath string) *Bastion {
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
func (app *Bastion) Serve() error {
	server := http.Server{Addr: app.cfg.Server.Addr, Handler: app.r}

	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	// check for a closing signal
	go func() {
		// graceful shutdown
		<-ctx.Done()
		log.Printf("shutting down application")

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("unable to shutdown server: %v", err)
		} else {
			log.Printf("server stopped")
		}
	}()

	// start the web server
	log.Printf("Starting application at %s\n", app.cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func initialize(app *Bastion) {
	/**
	 * internal router
	 */
	app.r = chi.NewRouter()
	app.r.Use(middleware.Recovery)

	/**
	 * Ping route
	 */
	app.r.Get("/ping", pingHandler)

	/**
	 * API Router
	 */
	app.APIRouter = chi.NewRouter()
	app.APIRouter.Use(CHIMiddleware.RequestID)
	app.APIRouter.Use(CHIMiddleware.Logger)
	app.r.Mount(app.cfg.API.BasePath, app.APIRouter)
}

// NewRouter return a router as a subrouter along a routing path.
// It's very useful to split up a large API as many independent routers and
// compose them as a single service.
func NewRouter() chi.Router {
	return chi.NewRouter()
}
