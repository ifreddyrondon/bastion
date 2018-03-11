package bastion

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi"
	CHIMiddleware "github.com/go-chi/chi/middleware"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// onShutdown is a function to be implemented when is necessary
// to run something before a shutdown of the server or in graceful shutdown.
type onShutdown func()

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
// Without a Bastion you can't do much!
type Bastion struct {
	r      *chi.Mux
	server *http.Server
	*Options
	APIRouter *chi.Mux
}

// New returns a new instance of Bastion and adds some sane, and useful, defaults.
// 	Defaults:
//		Addr: "127.0.0.1:8080"
//		Env: "development"
//		Debug: false
//		API:
//			BasePath: "/"
func New(opts Options) *Bastion {
	app := new(Bastion)
	app.Options = optionsWithDefaults(&opts)
	initialize(app)
	app.server = &http.Server{Addr: app.Options.Addr, Handler: app.r}
	return app
}

// FromFile is an util function to load  a new instance of Bastion from a options file.
// The options file could it be in YAML or JSON format. Is some attributes are missing
// from the config file it'll be set with the defaults.
// FromFile takes a special consideration for `server.address` default.
// When it's not provided it'll search the ADDR and PORT environment variables
// first before set the default.
func FromFile(path string) (*Bastion, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "missing configuration file at %v", path)
	}

	var opts Options
	if err := yaml.Unmarshal(b, &opts); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal configuration file")
	}
	app := new(Bastion)
	app.Options = optionsWithDefaults(&opts)
	initialize(app)
	app.server = &http.Server{Addr: app.Options.Addr, Handler: app.r}
	return app, nil
}

// RegisterOnShutdown registers a function to call on Shutdown.
// This can be used to gracefully shutdown connections that have
// undergone NPN/ALPN protocol upgrade or that have been hijacked.
// This function should start protocol-specific graceful shutdown,
// but should not wait for shutdown to complete.
func (app *Bastion) RegisterOnShutdown(fs ...onShutdown) {
	for _, f := range fs {
		app.server.RegisterOnShutdown(f)
	}
}

// Serve accepts incoming incoming connections coming from the specified address/port.
// It also prepare the graceful shutdown.
func (app *Bastion) Serve() error {
	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	go graceful(ctx, app.server)
	// start the web server
	log.Printf("[app:starting] at %s\n", app.Options.Addr)
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
	app.r.Use(Recovery)

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
	app.r.Mount(app.Options.APIBasepath, app.APIRouter)
}

// NewRouter return a router as a subrouter along a routing path.
// It's very useful to split up a large API as many independent routers and
// compose them as a single service.
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
