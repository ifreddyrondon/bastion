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
    res := struct {Message string `json:"message"`}{Message: "world"}
    render.NewJSON().Send(w, res)
}

func main() {
    app := bastion.New()
    app.Get("/hello", handler)
    app.Serve()
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

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
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
	w.Write([]byte("todos list of stuff.."))
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos create"))
}

func get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("get todo with id %v", id)))
}

func update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("update todo with id %v", id)))
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("delete todo with id %v", id)))
}

func main() {
	app := bastion.New()
	app.Mount("/todo/", routes())
	app.Serve()
}
```

## Middleware

Bastion comes equipped with a set of commons middleware handlers, providing a suite of standard `net/http` middleware.
They are just stdlib net/http middleware handlers. There is nothing special about them, which means the router and all 
the tooling is designed to be compatible and friendly with any middleware in the community.

### Core middleware

Name | Description
---- | -----------
Logger | Logs the start and end of each request with the elapsed processing time.
Recovery | Gracefully absorb panics and prints the stack trace.
RequestID | Injects a request ID into the context of each request.
APIErrHandler | Intercept responses to verify if his status code is >= 500. If status is >= 500, it'll response with a [default error](#api500errmessage). IT allows to response with the same error without disclosure internal information, also the real error is logged.

### Auxiliary middleware

Name | Description
---- | -----------
Listing | Parses the url from a request and stores a [listing.Listing](https://github.com/ifreddyrondon/bastion/blob/master/middleware/listing/listing.go#L11) on the context, it can be accessed through middleware.GetListing.

For more references check [chi middleware](https://github.com/go-chi/chi/tree/master#middlewares)

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
    app.Serve()
}
```

## Options

Options are used to define how the application should run.

### Structure

```go
// Options are used to define how the application should run.
type Options struct {
	// API500ErrMessage message returned to the user when catch a 500 status error.
	API500ErrMessage string `yaml:"api500ErrMessage"`
	// Addr bind address provided to http.Server. Default is "127.0.0.1:8080"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `yaml:"addr"`
	// Env "environment" in which the App is running. Default is "development".
	Env string `yaml:"env"`
	// NoPrettyLogging don't output a colored human readable version on the out writer.
	NoPrettyLogging bool `yaml:"noPrettyLogging"`
	// LoggerLevel defines log levels. Default is DebugLevel.
	LoggerLevel Level `yaml:"loggerLevel"`
	// LoggerOutput logger output writer. Default os.Stdout
	LoggerOutput io.Writer
}
```

#### `API500ErrMessage`

Api 500 error message represent the message returned to the user when a http 500 error is caught by the APIErrHandler middleware. Default `looks like something went wrong`. It's JSON tagged as `api500ErrMessage`

When:

```json
"api500ErrMessage": "looks like something went wrong",
```

Then: `http://localhost:8080/foo/test`

#### `Addr`

Addr is the bind address provided to http.Server. Default is `127.0.0.1:8080`. Can be set using **ENV** vars `ADDR` and `PORT`. It's JSON tagged as `addr`

#### Env

Env is the "environment" in which the App is running. Default is "development". Can be set using **ENV** vars `GO_ENV` It's JSON tagged as `env`

#### NoPrettyLogging

NoPrettyLogging boolean flag to don't output a colored human readable version on the out writer. Default `false`. It's JSON tagged as `noPrettyLogging`

#### LoggerLevel

LoggerLevel defines log levels. Allows for logging at the following levels (from highest to lowest):

- panic (`bastion.PanicLevel`, 5)
- fatal (`bastion.FatalLevel`, 4)
- error (`bastion.ErrorLevel`, 3)
- warn (`bastion.WarnLevel`, 2)
- info (`bastion.InfoLevel`, 1)
- debug (`bastion.DebugLevel`, 0)

Default `bastion.DebugLevel`, to turn off logging entirely, pass the bastion.Disabled constant. It's JSON tagged as `loggerLevel`.

#### LoggerOutput

LoggerOutput is an `io.Writer` where the logger output write. Default os.Stdout.

Each logging operation makes a single call to the Writer's Write method. There is no guaranty on access serialization to the Writer. If your Writer is not thread safe, you may consider using sync wrapper.

### From optionals functions

Bastion can be configured with optionals functions that are optional when using `bastion.New()`.

- `API500ErrMessage(msg string)` set the message returned to the user when catch a 500 status error.
- `Addr(add string)` bind address provided to http.Server.
- `Env(env string)` set the "environment" in which the App is running.
- `NoPrettyLogging()` turn off the pretty logging.
- `LoggerLevel(lvl Level)` set the logger level.
- `LoggerOutput(w io.Writer)` set the logger output writer.

```go
package main

import (
    "github.com/ifreddyrondon/bastion"
)

func main() {
	bastion.New() // defaults options
	bastion.New(bastion.NoPrettyLogging(), bastion.Addr("0.0.0.0:3000")) // turn off pretty print logger and sets address to 0.0.0.0:3000
}
```

### From options file

Bastion comes with an util function to load a new instance of Bastion from a options file. The options file could it be in **YAML** or **JSON** format. Is some attributes are missing from the options file it'll be set with the default. [Example](https://github.com/ifreddyrondon/bastion/blob/master/_examples/options-yaml/main.go).

FromFile takes special consideration when there are **ENV** vars:

* For `Addr`. When it's not provided it'll search the `ADDR` and `PORT` environment variables first before set the default.

* For `Env`. When it's not provided it'll search the `GO_ENV` environment variables first before set the default.

#### YAML

```yaml
addr: ":8080"
debug: true
env: "development"
```

#### JSON

```json
{
  "addr": ":8080",
  "debug": true,
  "env": "development"
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

Render a HTTP status code and content type to the associated Response.

### StringRenderer
- **render.Text** response strings with text/plain Content-Type.
```go
render.Text.Response(rr, http.StatusOK, "test")
```
- **render.HTML** response strings with text/html Content-Type.
```go
render.HTML.Response(rr, http.StatusOK, "<h1>Hello World</h1>")
```

### ByteRenderer
- **render.Data** response []byte with application/octet-stream Content-Type.
```go
render.Data.Response(rr, http.StatusOK, []byte("test"))
```

### Renderer

Handle the marshaler of structs responses to the client.

```go
// Renderer interface for managing response payloads.
type Renderer interface {
	// Response encoded responses in the ResponseWriter with the HTTP status code.
	Response(w http.ResponseWriter, code int, response interface{})
}
```

APIRenderer are convenient methods for api responses.

```go

// APIRenderer interface for managing API response payloads.
type APIRenderer interface {
	Renderer
	OKRenderer
	ClientErrRenderer
	ServerErrRenderer
}

// OKRenderer interface for managing success API response payloads.
type OKRenderer interface {
	Send(w http.ResponseWriter, response interface{})
	Created(w http.ResponseWriter, response interface{})
	NoContent(w http.ResponseWriter)
}

// ClientErrRenderer interface for managing API responses when client error.
type ClientErrRenderer interface {
	BadRequest(w http.ResponseWriter, err error)
	NotFound(w http.ResponseWriter, err error)
	MethodNotAllowed(w http.ResponseWriter, err error)
}

// ServerErrRenderer interface for managing API responses when server error.
type ServerErrRenderer interface {
	InternalServerError(w http.ResponseWriter, err error)
}
```

[JSON](https://github.com/ifreddyrondon/bastion/blob/master/render/json.go) and [XML](https://github.com/ifreddyrondon/bastion/blob/master/render/xml.go) implements APIRenderer and they can be configured with optional functions.

#### E.g.

- [JSON](https://github.com/ifreddyrondon/bastion/blob/master/render/json_test.go)
- [XML](https://github.com/ifreddyrondon/bastion/blob/master/render/xml_test.go)

Response a JSON with a 200 HTTP status code.

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{Message: "world"}
	render.NewJSON().Send(w, res)
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	app.Serve()
}
```

## Logger

Bastion commes with a JSON structured logger powered by [github.com/rs/zerolog](github.com/rs/zerolog). It can be accessed through the bastion instance `bastion.Logger` or from the context of each request `l := bastion.LoggerFromCtx(ctx)`

### Logging from bastion instance

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

func main() {
	app := bastion.New()
	app.Logger.Info().Str("app", "demo").Msg("main")
	app.Serve()
}
```

### Logging from handler

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

func handler(w http.ResponseWriter, r *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{Message: "world"}
	l := bastion.LoggerFromCtx(r.Context())
	l.Info().Msg("handler")

	render.NewJSON().Send(w, res)
}

func main() {
	app := bastion.New()
	app.Get("/hello", handler)
	app.Serve()
}
```
