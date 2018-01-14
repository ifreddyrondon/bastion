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
	"github.com/markbates/sigtx"
)

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	r          *chi.Mux
	cfg        *config.Config
	APIRouter  *chi.Mux
	finalizers []Finalizer
	Reader
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
	app.Reader = new(JsonReader)
	app.Responder = new(JsonResponder)
	initialize(app)
	return app
}

// AppendFinalizers add helpers function that will be executed
// in the graceful shutdown.
// The function need to implement the Finalizer interface.
func (app *Bastion) AppendFinalizers(finalizer ...Finalizer) {
	app.finalizers = append(app.finalizers, finalizer...)
}

// Serve the application at the specified address/port
func (app *Bastion) Serve() error {
	server := http.Server{Addr: app.cfg.Server.Addr, Handler: app.r}

	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	appFinalizer := serverFinalizer{&server, ctx}
	app.AppendFinalizers(appFinalizer)

	// check for a closing signal
	go func() {
		// graceful shutdown
		<-ctx.Done()

		log.Printf("shutting down application")
		if err := finalize(app.finalizers); err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("application stopped")
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
