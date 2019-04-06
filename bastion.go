package bastion

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/ifreddyrondon/bastion/middleware"
)

// onShutdown is a function to be implemented when is necessary
// to run something before a shutdown of the server or in graceful shutdown.
type onShutdown func()

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
// Without a Bastion you can't do much!
type Bastion struct {
	server *http.Server
	Logger *zerolog.Logger
	Options
	*chi.Mux
}

// New returns a new instance of Bastion and adds some sane, and useful, defaults.
// 	Defaults:
//		Addr: "127.0.0.1:8080"
//		Env: "development"
//		Debug: false
func New(opts ...Opt) *Bastion {
	app := &Bastion{}
	for _, opt := range opts {
		opt(app)
	}
	setDefaultsOpts(&app.Options)

	initialize(app)
	return app
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

	go graceful(ctx, app)

	app.Logger.Info().Msgf("app starting at %v", app.Options.Addr)
	if err := app.server.ListenAndServe(); err != nil {
		app.Logger.Error().Err(err).Msg("listenAndServe err")
		return err
	}
	return nil
}

func initialize(app *Bastion) {
	/**
	 * init logger
	 */
	app.Logger = getLogger(&app.Options)

	/**
	 * router
	 */
	app.Mux = chi.NewRouter()
	app.Mux.Use(hlog.NewHandler(*app.Logger))
	apiErr := middleware.APIError(
		middleware.APIErrorDefault500(errors.New(app.Options.API500ErrMessage)),
		middleware.APIErrorLoggerOutput(app.Options.LoggerOutput),
	)
	app.Mux.Use(apiErr)
	recovery := middleware.Recovery(middleware.RecoveryLoggerOutput(app.Options.LoggerOutput))
	app.Mux.Use(recovery)
	app.Mux.Use(loggerRequest(!app.Options.isDEV())...)

	/**
	 * Ping route
	 */
	app.Mux.Get("/ping", pingHandler)
	app.server = &http.Server{Addr: app.Options.Addr, Handler: app.Mux}
}

// NewRouter return a router as a subrouter along a routing path.
// It's very useful to split up a large API as many independent routers and
// compose them as a single service.
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
