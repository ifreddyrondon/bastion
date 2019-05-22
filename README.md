# Bastion

[![Documentation](https://godoc.org/github.com/ifreddyrondon/bastion?status.svg)](http://godoc.org/github.com/ifreddyrondon/bastion)
[![Coverage Status](https://coveralls.io/repos/github/ifreddyrondon/bastion/badge.svg?branch=master)](https://coveralls.io/github/ifreddyrondon/bastion?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ifreddyrondon/bastion)](https://goreportcard.com/report/github.com/ifreddyrondon/bastion)
[![CircleCI](https://circleci.com/gh/ifreddyrondon/bastion.svg?style=svg)](https://circleci.com/gh/ifreddyrondon/bastion)

Defend your API from the sieges. Bastion offers an "augmented" Router instance.

It has the minimal necessary to create an API with default handlers and middleware that help you raise your API easy and fast.
Allows to have commons handlers and middleware between projects with the need for each one to do so. It's also included some 
useful/optional subpackages: [middleware](https://github.com/ifreddyrondon/bastion/blob/master/middleware) and [render](https://github.com/ifreddyrondon/bastion/blob/master/render). We hope you enjoy it too!

## Installation

`go get -u github.com/ifreddyrondon/bastion`

## Examples

See [_examples/](https://github.com/ifreddyrondon/bastion/blob/master/_examples/) for a variety of examples.

**As easy as:**

```go
package main

import (
    "net/http"

    "github.com/ifreddyrondon/bastion"
    "github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, r *http.Request) {
	render.JSON.Send(w, map[string]string{"message": "hello bastion"})
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	// By default it serves on :8080 unless a
	// ADDR environment variable was defined.
	app.Serve()
	// app.Serve(":3000") for a hard coded port
}
```

## Router

Bastion use [go-chi](https://github.com/go-chi/chi) as a router making it easy to modularize the applications. 
Each Bastion instance accepts a URL `pattern` and chain of `handlers`. The URL pattern supports 
named params (ie. `/users/{userID}`) and wildcards (ie. `/admin/*`). URL parameters can be fetched 
at runtime by calling `chi.URLParam(r, "userID")` for named parameters and `chi.URLParam(r, "*")` 
for a wildcard parameter.

### NewRouter

NewRouter return a router as a subrouter along a routing path.

It's very useful to split up a large API as many independent routers and compose them as a single service.

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Routes creates a REST router for the todos resource
func routes() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", list)    // GET /todos - read a list of todos
	r.Post("/", create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func list(w http.ResponseWriter, r *http.Request) {
	render.Text.Response(w, http.StatusOK, "todos list of stuff..")
}

func create(w http.ResponseWriter, r *http.Request) {
	render.Text.Response(w, http.StatusOK, "todos create")
}

func get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	render.Text.Response(w, http.StatusOK, fmt.Sprintf("get todo with id %v", id))
}

func update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	render.Text.Response(w, http.StatusOK, fmt.Sprintf("update todo with id %v", id))
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	render.Text.Response(w, http.StatusOK, fmt.Sprintf("delete todo with id %v", id))
}

func main() {
	app := bastion.New()
	app.Mount("/todo/", routes())
	fmt.Fprintln(os.Stderr, app.Serve())
}
```

## Middlewares

Bastion comes equipped with a set of commons middleware handlers, providing a suite of standard `net/http` middleware.
They are just stdlib net/http middleware handlers. There is nothing special about them, which means the router and all 
the tooling is designed to be compatible and friendly with any middleware in the community.

### Core middleware

Name | Description
---- | -----------
Logger | Logs the start and end of each request with the elapsed processing time.
RequestID | Injects a request ID into the context of each request.
Recovery | Gracefully absorb panics and prints the stack trace.
InternalError | Intercept responses to verify if his status code is >= 500. If status is >= 500, it'll response with a [default error](#InternalErrMsg). It allows to response with the same error without disclosure internal information, also log real error (default callback implementation. Check `InternalErrCallback` func to override it).

### Auxiliary middleware

Name | Description
---- | -----------
Listing | Parses the url from a request and stores a [listing.Listing](https://github.com/ifreddyrondon/bastion/blob/master/middleware/listing/listing.go#L11) on the context, it can be accessed through middleware.GetListing.
WrapResponseWriter | provides an easy way to capture http related metrics from your application's http.Handlers or event hijack the response. 

Checkout for references, examples, options and docu in [middleware](https://github.com/ifreddyrondon/bastion/blob/master/middleware) or [chi](https://github.com/go-chi/chi/tree/master#middlewares) for more middlewares. 

## Register on shutdown

You can register a function to call on shutdown. This can be used to gracefully shutdown connections. By default the shutdown execute the server shutdown.

Bastion listens if any **SIGINT**, **SIGTERM** or **SIGKILL** signal is emitted and performs a graceful shutdown.

It can be added with `RegisterOnShutdown` method of the bastion instance, it can accept variable number of functions.

### Register on shutdown example

```go
package main

import (
    "log"

    "github.com/ifreddyrondon/bastion"
)

func onShutdown() {
    log.Printf("My registered on shutdown. Doing something...")
}

func main() {
    app := bastion.New()
    app.RegisterOnShutdown(onShutdown)
    app.Serve(":8080")
}
```

## Options

Options are used to define how the application should run, it can be set through optionals functions when using `bastion.New()`.

```go
package main

import (
    "github.com/ifreddyrondon/bastion"
)

func main() {
	// turn off pretty print logger and sets 500 errors message
	bastion.New(bastion.DisablePrettyLogging(), bastion.InternalErrMsg(`Just another "500 - internal error"`))
}
```

### InternalErrMsg

Represent the message returned to the user when a http 500 error is caught by the InternalError middleware. 
Default `looks like something went wrong`.

- `InternalErrMsg(msg string)` set the message returned to the user when catch a 500 status error.

### InternalErrCallback

Callback function to handler the real error catched by InternalError middleware.

- `InternalErrCallback(f func(int, io.Reader))` sets the callback function when internal error middleware catch a 500 error.

### DisableInternalErrorMiddleware

Boolean flag to disable the [internal error middleware](https://github.com/go-chi/chi/tree/master#middlewares). Default `false`.

- `DisableInternalErrorMiddleware()` turn off internal error middleware.

### DisableRecoveryMiddleware

Boolean flag to disable [recovery middleware](https://github.com/go-chi/chi/tree/master#middlewares). Default `false`.

- `DisableRecoveryMiddleware()` turn off recovery middleware.

### DisablePingRouter

Boolean flag to disable the ping route. Default `false`.

- `DisablePingRouter()` turn off ping route.

### DisableLoggerMiddleware

Boolean flag to disable the logger middleware. Default `false`.

- `DisableLoggerMiddleware()` turn off logger middleware.

### DisablePrettyLogging

Boolean flag to don't output a colored human readable version on the out writer. Default `false`.

- `DisablePrettyLogging()` turn off the pretty logging.

### LoggerLevel

Defines log level. Default `debug`. Allows for logging at the following levels (from highest to lowest):

- panic, 5
- fatal, 4
- error, 3
- warn, 2
- info, 1
- debug, 0

- `LoggerLevel(lvl string)` set the logger level.

```go
package main

import (
    "github.com/ifreddyrondon/bastion"
)

func main() {
	bastion.New(bastion.LoggerLevel(bastion.ErrorLevel))
	// or
	bastion.New(bastion.LoggerLevel("error"))
}
```

### LoggerOutput

Where the logger output write. Default `os.Stdout`.

- `LoggerOutput(w io.Writer)` set the logger output writer.

### ProfilerRoutePrefix 

Optional path prefix for profiler subrouter. If left unspecified, `/debug/` is used as the default path prefix.

- `ProfilerRoutePrefix(prefix string)` set the prefix path for the profile router.

### DisableProfiler 

Boolean flag to disable the profiler subrouter.

- `DisableProfiler()` turn off profiler subrouter.

### IsProduction()

IsProduction check if app is running in production mode.
 
Can be set using `ProductionMode()` option or with **ENV** vars `GO_ENV` or `GO_ENVIRONMENT`. 
`ProductionMode()` has more priority than the ENV variables. 

When **production** mode is on, the request logger IP, UserAgent and Referer are enable, the logger level is set 
to `error`, the profiler routes are disabled and the logging pretty print is disabled.

- `ProductionMode()` set the app to production mode.
- `ProductionMode(false)` or force debug.

```go
package main

import (
    "github.com/ifreddyrondon/bastion"
)

func main() {
	bastion.New(bastion.ProductionMode())
	// or with bool param to force debug mode even with env vars.
	bastion.New(bastion.ProductionMode(false))
}
```

## Testing

Bastion comes with battery included testing tools to perform End-to-end test over your endpoint/handlers.

It uses [github.com/gavv/httpexpect](https://github.com/gavv/httpexpect) to incrementally build HTTP requests,
inspect HTTP responses and inspect response payload recursively.

### Quick start

1. Create the bastion instance with the handler you want to test.
2. Import from `bastion.Tester`
3. It receive a `*testing.T` and `*bastion.Bastion` instances as params.
4. Build http request.
5. Inspect http response.
6. Inspect response payload.

```go
package main_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
    "github.com/ifreddyrondon/bastion/render"
)

func setup() *bastion.Bastion {
	app := bastion.New()
	app.Mount("/todo/", todo.Routes())
	return app
}

func TestHandlerCreate(t *testing.T) {
	app := setup()
	payload := map[string]interface{}{
		"description": "new description",
	}

	e := bastion.Tester(t, app)
	e.POST("/todo/").WithJSON(payload).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 0).
		ContainsKey("description").ValueEqual("description", "new description")
}
```

Go and check the [full test](https://github.com/ifreddyrondon/bastion/blob/master/_examples/todo-rest/todo/handler_test.go) for [handler](https://github.com/ifreddyrondon/bastion/blob/master/_examples/todo-rest/todo/handler.go) and complete [app](https://github.com/ifreddyrondon/bastion/tree/master/_examples/todo-rest) 🤓

## Render

Easily rendering JSON, XML, binary data, and HTML templates responses 

### Usage
It can be used with pretty much any web framework providing you can access the `http.ResponseWriter` from your handler.
The rendering functions simply wraps Go's existing functionality for marshaling and rendering data.

```go
package main

import (
	"encoding/xml"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func main() {
	app := bastion.New()

	app.Get("/data", func(w http.ResponseWriter, req *http.Request) {
		render.Data.Response(w, http.StatusOK, []byte("Some binary data here."))
	})

	app.Get("/text", func(w http.ResponseWriter, req *http.Request) {
		render.Text.Response(w, http.StatusOK, "Plain text here")
	})

	app.Get("/html", func(w http.ResponseWriter, req *http.Request) {
		render.HTML.Response(w, http.StatusOK, "<h1>Hello World</h1>")
	})

	app.Get("/json", func(w http.ResponseWriter, req *http.Request) {
		render.JSON.Response(w, http.StatusOK, map[string]string{"hello": "json"})
	})

	app.Get("/json-ok", func(w http.ResponseWriter, req *http.Request) {
		// with implicit status 200
		render.JSON.Send(w, map[string]string{"hello": "json"})
	})

	app.Get("/xml", func(w http.ResponseWriter, req *http.Request) {
		render.XML.Response(w, http.StatusOK, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Get("/xml-ok", func(w http.ResponseWriter, req *http.Request) {
		// with implicit status 200
		render.XML.Send(w, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Serve()
}
```

Checkout more references, examples, options and implementations in [render](https://github.com/ifreddyrondon/bastion/blob/master/render).

## Binder

To bind a request body or a source input into a type, use a binder. It's currently support binding of JSON, XML and YAML.
The binding execute `Validate()` if the type implements the `binder.Validate` interface after successfully bind the type.

The goal of implement `Validate` is to endorse the values linked to the type. This library intends for you to handle 
your own validations error.

### Usage

```go
package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"
)

type address struct {
	Address *string `json:"address" xml:"address" yaml:"address"`
	Lat     float64 `json:"lat" xml:"lat" yaml:"lat"`
	Lng     float64 `json:"lng" xml:"lng" yaml:"lng"`
}

func (a *address) Validate() error {
	if a.Address == nil || *a.Address == "" {
		return errors.New("missing address field")
	}
	return nil
}

func main() {
	app := bastion.New()
	app.Post("/decode-json", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.JSON.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Post("/decode-xml", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.XML.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Post("/decode-yaml", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.YAML.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Serve()
}
```

Checkout more references, examples, options and implementations in [binder](https://github.com/ifreddyrondon/bastion/blob/master/binder).

## Logger

Bastion have an internal JSON structured logger powered by [github.com/rs/zerolog](github.com/rs/zerolog). 
It can be accessed from the context of each request `l := bastion.LoggerFromCtx(ctx)`. The request id is logged for 
every call to the logger.

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, r *http.Request) {
	l := bastion.LoggerFromCtx(r.Context())
	l.Info().Msg("handler")

	render.JSON.Send(w, map[string]string{"message": "hello bastion"})
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	app.Serve()
}
```
