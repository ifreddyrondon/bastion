package bastion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/markbates/sigtx"
	"github.com/rs/zerolog"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/render"
)

const defaultAddr = ":8080"

// OnShutdown is a function to be implemented when is necessary
// to run something before a shutdown of the server or in graceful shutdown.
type OnShutdown func()

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
// Without a Bastion you can't do much!
type Bastion struct {
	r      *chi.Mux
	server *http.Server
	logger zerolog.Logger
	Options
	*chi.Mux
}

// New returns a new instance of Bastion and adds some sane, and useful, defaults.
func New(opts ...Opt) *Bastion {
	app := &Bastion{
		server: &http.Server{},
		Options: Options{
			InternalErrMsg:      defaultInternalErrMsg,
			ProfilerRoutePrefix: defaultProfilerRoutePrefix,
			LoggerLevel:         DebugLevel,
			LoggerOutput:        os.Stdout,
		},
	}
	for _, opt := range opts {
		opt(app)
	}
	app.codeMode = resolveMode(app.Mode)
	if !app.IsDebug() {
		app.DisablePrettyLogging = true
		app.DisableProfiler = true
		app.LoggerLevel = ErrorLevel
	}

	app.Mode = app.codeMode.String()

	l, err := getLogger(app.LoggerOutput, !app.DisablePrettyLogging, app.LoggerLevel)
	if err != nil {
		panic(err)
	}
	app.r = router(app.Options, *l)
	app.Mux = chi.NewMux()
	app.r.Mount("/", app.Mux)
	app.logger = l.With().Str("module", "bastion").Logger()

	if app.IsDebug() {
		app.logger.Info().Msg(`Running in "debug" mode. Switch to "production" mode in production.
 - using code:  bastion.New(bastion.Mode("production"))
 - using env: export GO_ENV=production
 - using env: export GO_ENVIRONMENT=production

`)
	}

	return app
}

func router(opts Options, l zerolog.Logger) *chi.Mux {
	mux := chi.NewMux()
	mux.NotFound(notFound)
	mux.MethodNotAllowed(notAllowed)
	// logger middleware
	if !opts.DisableLoggerMiddleware {
		logMiddleware := []middleware.LoggerOpt{
			middleware.AttachLogger(l),
		}
		if !opts.IsDebug() {
			logMiddleware = append(
				logMiddleware,
				middleware.EnableLogReferer(),
				middleware.EnableLogUserAgent(),
				middleware.EnableLogReqIP(),
			)
		}
		logger := middleware.Logger(logMiddleware...)
		mux.Use(logger)
	}

	// internal error middleware
	if !opts.DisableInternalErrorMiddleware {
		internalErr := middleware.InternalError(
			middleware.InternalErrMsg(errors.New(opts.InternalErrMsg)),
			middleware.InternalErrLoggerOutput(opts.LoggerOutput),
		)
		mux.Use(internalErr)
	}

	// recovery middleware
	if !opts.DisableRecoveryMiddleware {
		recovery := middleware.Recovery(middleware.RecoveryLoggerOutput(opts.LoggerOutput))
		mux.Use(recovery)
	}

	if !opts.DisablePingRouter {
		mux.Get("/ping", pingHandler)
	}
	if !opts.DisableProfiler {
		mux.Mount(opts.ProfilerRoutePrefix, chiMiddleware.Profiler())
	}

	return mux
}

func notFound(w http.ResponseWriter, r *http.Request) {
	render.JSON.NotFound(w, fmt.Errorf("resource %s not found", r.URL.Path))
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("method %s not allowed for resource %s", r.Method, r.URL.Path)
	render.JSON.MethodNotAllowed(w, err)
}

func printRoutes(mux *chi.Mux, opts Options, l *zerolog.Logger) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.HasPrefix(route, opts.ProfilerRoutePrefix) {
			return nil
		}
		route = strings.Replace(route, "/*/", "/", -1)
		l.Info().Str("component", "route").Msgf("%s %s", method, route)
		return nil
	}

	if err := chi.Walk(mux, walkFunc); err != nil {
		l.Error().Err(err).Msgf("walking through the routes")
	}
}

// RegisterOnShutdown registers a function to call on Shutdown.
// This can be used to gracefully shutdown connections that have
// undergone NPN/ALPN protocol upgrade or that have been hijacked.
// This function should start protocol-specific graceful shutdown,
// but should not wait for shutdown to complete.
func (app *Bastion) RegisterOnShutdown(fs ...OnShutdown) {
	for _, f := range fs {
		app.server.RegisterOnShutdown(f)
	}
}

// Serve accepts incoming connections coming from the specified address/port.
// It is a shortcut for http.ListenAndServe(addr, router).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (app *Bastion) Serve(addr ...string) error {
	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	go graceful(ctx, app.server, &app.logger)

	address := resolveAddress(addr, &app.logger)
	app.logger.Info().Msgf("app starting at %v", address)
	app.server.Addr = address
	app.server.Handler = app.r

	printRoutes(app.r, app.Options, &app.logger)
	if err := app.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			app.logger.Info().Str("component", "Serve").Msg("http: Server closed")
			return err
		}
		app.logger.Error().Str("component", "Serve").Err(err).Msg("listenAndServe")
		return err
	}
	return nil
}

func graceful(ctx context.Context, server *http.Server, l *zerolog.Logger) {
	<-ctx.Done()
	logger := l.With().Str("component", "graceful").Logger()
	logger.Info().Msg("preparing for shutdown")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err)
		return
	}
	logger.Info().Msg("gracefully stopped")
}

func resolveAddress(addr []string, l *zerolog.Logger) string {
	switch len(addr) {
	case 0:
		if envAddr := os.Getenv("ADDR"); envAddr != "" {
			l.Debug().Msgf(`Environment variable ADDR="%s"`, envAddr)
			return envAddr
		}
		l.Debug().Msg("Environment variable ADDR is undefined. Using addr :8080 by default")
		return defaultAddr
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

// NewRouter return a router as a subrouter along a routing path.
// It's very useful to split up a large API as many independent routers and
// compose them as a single service.
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
