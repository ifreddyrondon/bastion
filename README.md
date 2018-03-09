# Bastion

Defend your API from the sieges. Bastion offers an "augmented" Router instance.

It has the minimal necessary to create an API with default handlers and middleware that help you raise your API easy and fast.
Allows to have commons handlers and middleware between projects with the need for each one to do so.

## Examples

* [helloworld](https://github.com/ifreddyrondon/bastion/blob/master/_examples/helloworld/main.go) - Quickstart, first Hello world with bastion.
* [todo-rest](https://github.com/ifreddyrondon/bastion/blob/master/_examples/todo-rest/) - REST APIs made easy, productive and maintainable.
* [Config with yaml](https://github.com/ifreddyrondon/bastion/blob/master/_examples/config-yaml/main.go) - Bastion with config file.
* [Register on shutdown](https://github.com/ifreddyrondon/bastion/blob/master/_examples/register/main.go) - Registers functions to be call on Shutdown.

## Table of contents

* [Installation](#installation)
* [Router](#router)
	* [NewRouter](#newrouter)
	* [Example](#example)
* [Middlewares](#middlewares)
* [Register on shutdown](#register-on-shutdown)
	* [Example](#example-1)
* [Configuration](#configuration)
	* [Structure](#structure)
		* [Api](#api)
		* [Server](#server)
		* [Debug](#debug)
	* [From configuration file](#from-configuration-file)
		* [YAML](#yaml)
		* [JSON](#json)
* [Testing](#testing)
	* [Quick start](#quick-start)
* [Render](#render)
    * [Example](#example-2)

## Installation

`go get -u github.com/ifreddyrondon/bastion`

## Router
Bastion use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
of mounting an api router, it will define the routes and middleware of the application with the app logic.

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

type Handler struct{}

// Routes creates a REST router for the todos resource
func (h *Handler) Routes() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", h.list)    // GET /todos - read a list of todos
	r.Post("/", h.create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", h.update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", h.delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos create"))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("get todo with id %v", id)))
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("update todo with id %v", id)))
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("delete todo with id %v", id)))
}

func main() {
	app := bastion.New(nil)
	app.APIRouter.Mount("/todo/", new(Handler).Routes())
	app.Serve()
}
```

### Example

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {Message string `json:"message"`}{"world"}
	json.NewRender(w).Send(res)
}

func main() {
	app := bastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
```

## Middlewares

Bastion comes equipped with a set of commons middlewares, providing a suite of standard
`net/http` middlewares.

Name | Description
---- | -----------
Logger | Logs the start and end of each request with the elapsed processing time
Recovery | Gracefully absorb panics and prints the stack trace
RequestID | Injects a request ID into the context of each request

## Register on shutdown
You can register a function to call on shutdown. This can be used to gracefully shutdown connections. By default the shutdown execute the server shutdown.

Bastion listens if any **SIGINT**, **SIGTERM** or **SIGKILL** signal is emitted and performs a graceful shutdown.

It can be added with `RegisterOnShutdown` method of the bastion instance, it can accept variable number of functions.

### Example

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
	app := bastion.New(nil)
	app.RegisterOnShutdown(onShutdown)
	app.Serve()
}
```

## Configuration
Represents the configuration for bastion. Config are used to define how the application should run.

### Structure
```go
type Config struct {
 	API struct {
 		BasePath string
 	}
 	Server struct {
 		Addr string
 	}
 	Debug bool
}
```

#### Api
##### `Api.BasePath`
Base path value where the application is going to be mounted. Default `/`. Is JSON tagged as `api.base_path`

When
```json
"base_path": "/foo/test",
```
Then
```
http://localhost:8080/foo/test
```

#### Server
##### `Server.Addr`
Address is the host and port where the app is serve. Default `127.0.0.1:8080`. Is JSON tagged as `server.address`

#### Debug
Debug flag if Bastion should enable debugging features. Default `false`. . Is JSON tagged as `debug`

### From configuration file
Bastion comes with an util function to load configuration from a file.
**FromFile** is an util function to load the bastion configuration from a config file. The config file could it be in **YAML** or **JSON** format. Is some attributes are missing
from the config file it'll be set with the default. [Example](https://github.com/ifreddyrondon/bastion/blob/master/_examples/config-yaml/main.go).

FromFile takes a special consideration for `server.address` default. When it's not provided it'll search the ADDR and PORT environment variables first before set the default.

#### YAML
```yaml
api:
  base_path: "/"
server:
  address: ":8080"
debug: true

```
#### JSON
```json
{
  "api": {
    "base_path": "/"
  },
  "server": {
    "address": ":8080"
  },
  "debug": true
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
    "github.com/ifreddyrondon/bastion/render/json"
)

func setup() *bastion.Bastion {
	app := bastion.New(nil)
	handler := todo.Handler{
		Render: json.NewRender,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
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
The render engine implements `Engine` and is obtained through `Render` function.  

```go
// Engine define methods to encoded response in the body of a request with the HTTP status code.
type Engine interface {
	Response(code int, response interface{})
	Send(response interface{})
	Created(response interface{})
	NoContent()
	BadRequest(err error)
	NotFound(err error)
	MethodNotAllowed(err error)
	InternalServerError(err error)
}

// Render returns a Engine to response a request with the HTTP status code.
type Render func(http.ResponseWriter) Engine
``` 

Bastion define a `json.Render` [implementation](https://github.com/ifreddyrondon/bastion/blob/master/render/json/json.go) of `Engine` and is available through `json.NewRender`

### Example

Response a JSON with a 200 HTTP status code. 
```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	json.NewRender(w).Send(res)
}

func main() {
	app := bastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
```
