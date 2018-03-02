package bastion

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi"
	CHIMiddleware "github.com/go-chi/chi/middleware"
	"github.com/ifreddyrondon/bastion/config"
	"github.com/markbates/sigtx"
)

// DefaultResponder is the default Responder and is used to provide utils method
// for response http request by bastion. It's a JsonResponder instance and it'll
// wrap the response to a JSON valid response.
var DefaultResponder Responder = new(JsonResponder)

// onShutdown is a function to be implemented when is necessary
// to run something before a shutdown of the server or in graceful shutdown.
type onShutdown func()

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	r         *chi.Mux
	cfg       *config.Config
	server    *http.Server
	APIRouter *chi.Mux
	Responder
}

// New returns a new Bastion instance.
// if cfg is empty the configuration will be from defaults.
// 	Defaults:
//		Api:
//			BasePath: "/"
//		Server:
//			Addr ":8080"
//		Debug: false
func New(cfg *config.Config) *Bastion {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	app := new(Bastion)
	app.cfg = cfg
	app.Responder = DefaultResponder
	initialize(app)
	app.server = &http.Server{Addr: app.cfg.Server.Addr, Handler: app.r}
	return app
}

// onShutdown registers a function to call on Shutdown.
// This can be used to gracefully shutdown connections that have
// undergone NPN/ALPN protocol upgrade or that have been hijacked.
// This function should start protocol-specific graceful shutdown,
// but should not wait for shutdown to complete.
func (app *Bastion) RegisterOnShutdown(fs ...onShutdown) {
	for _, f := range fs {
		app.server.RegisterOnShutdown(f)
	}
}

// Serve, serve all the incoming connections coming from the specified address/port.
// It also prepare the graceful shutdown.
func (app *Bastion) Serve() error {
	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	go graceful(app.server, ctx)
	// start the web server
	log.Printf("[app:starting] at %s\n", app.cfg.Server.Addr)
	if err := app.server.ListenAndServe(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func initialize(app *Bastion) {
	/**
	 * internal router
	 */
	app.r = chi.NewRouter()
	app.r.Use(Recovery(app.Responder))

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
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
