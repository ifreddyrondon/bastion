# Bastion

Defend your API from the sieges. Bastion offers an "augmented" Router instance.

It has the minimal necessary to create an API with default handlers and middleware that help you raise your API easy and fast.
Allows to have commons handlers and middleware between projects with the need for each one to do so.

## Install

`go get -u github.com/ifreddyrondon/gobastion`

## Examples

* [helloworld](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/helloworld/main.go) - Quickstart, first Hello world with bastion.
* [todos-rest](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/todo-rest/) - REST APIs made easy, productive and maintainable.
* [config-yaml](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/config-yaml/main.go) - Bastion with config file.
* [finalizer](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/finalizer/main.go) - Bastion with Finalizer.

## Router
Bastion use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
of mounting an api router, it will define the routes and middleware of the application with the app logic.

### Example

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
)

var app *gobastion.Bastion

func handler(w http.ResponseWriter, _ *http.Request) {
	res := struct {Message string `json:"message"`}{"world"}
	app.Send(w, res)
}

func main() {
	app = gobastion.New(nil)
	app.APIRouter.Get("/hello", handler)
	app.Serve()
}
```

### NewRouter
NewRouter return a router as a subrouter along a routing path. 

It's very useful to split up a large API as many independent routers and compose them as a single service.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
)

type handler struct{}

// Routes creates a REST router for the todos resource
func (h *handler) Routes() chi.Router {
	r := gobastion.NewRouter()

	r.Get("/", h.List)    // GET /todos - read a list of todos
	r.Post("/", h.Create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", h.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", h.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos create"))
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("get todo with id %v", id)))
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("update todo with id %v", id)))
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("delete todo with id %v", id)))
}

func main() {
	app := gobastion.New(nil)
	app.APIRouter.Mount("/todo/", new(handler).Routes())
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

## Finalizers
Bastion listens if any **SIGINT**, **SIGTERM** or **SIGKILL** signal is emitted and performs a graceful shutdown.
By default the graceful shutdown execute the server shutdown through a Finalizer.

The Finalizer is an interface. All the finalizer will be executed into the graceful shutdown.
```go
type Finalizer interface {
	Finalize() error
}
```
It can be added to Finalizer queue with `AppendFinalizers` method of the bastion instance.

### Example

```go
package main

import (
	"log"

	"github.com/ifreddyrondon/gobastion"
)

type MyFinalizer struct{}

func (f MyFinalizer) Finalize() error {
	log.Printf("[finalizer:MyFinalizer] doing something")
	return nil
}

func main() {
	bastion := gobastion.New(nil)
	bastion.AppendFinalizers(MyFinalizer{})
	bastion.Serve()
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
from the config file it'll be set with the default. [Example](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/config-yaml/main.go).

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
2. Import from `gobastion.Tester`
3. It receive a `*testing.T` and `*gobastion.Bastion` instances as params.
4. Build http request.
5. Inspect http response.
6. Inspect response payload.

```go
var bastion *gobastion.Bastion

func TestMain(m *testing.M) {
	bastion = gobastion.New(nil)
	reader := new(gobastion.JsonReader)
	handler := todo.Handler{
		Reader:    reader,
		Responder: gobastion.DefaultResponder,
	}
	bastion.APIRouter.Mount("/todo/", handler.Routes())
	code := m.Run()
	os.Exit(code)
}

func TestHandlerCreate(t *testing.T) {
	payload := map[string]interface{}{
		"description": "new description",
	}

	e := gobastion.Tester(t, bastion)
	e.POST("/todo/").WithJSON(payload).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 0).
		ContainsKey("description").ValueEqual("description", "new description")
}
```

Go and check the [full test](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/todo-rest/todo/handler_test.go) for [handler](https://github.com/ifreddyrondon/gobastion/blob/master/_examples/todo-rest/todo/handler.go) and complete [app](https://github.com/ifreddyrondon/gobastion/tree/master/_examples/todo-rest) 🤓
