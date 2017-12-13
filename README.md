# Bastion

Defend your API from the sieges. Bastion offers an "augmented" Router instance.

It has the minimal necessary to create an API with default handlers and middleware that help you raise your API easy and fast.
Allows to have commons handlers and middleware between projects with the need for each one to do so.

## Examples

* [helloworld](https://github.com/ifreddyrondon/gobastion/blob/master/example/helloworld/main.go) - Quickstart, first Hello world with bastion.
* [config-yaml](https://github.com/ifreddyrondon/gobastion/blob/master/example/config-yaml/main.go) - Bastion with config file.
* [finalizer](https://github.com/ifreddyrondon/gobastion/blob/master/example/finalizer/main.go) - Bastion with Finalizer.
* [todos-rest](https://github.com/ifreddyrondon/gobastion/blob/master/example/todo-rest/main.go) - REST APIs made easy, productive and maintainable.

## Router
Bastion use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
of mounting an api router, it will define the routes and middleware of the application with the app logic.

### Example

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	utils.Send(w, res)
}

func main() {
	bastion := gobastion.New("")
	bastion.APIRouter.Get("/hello", helloHandler)
	bastion.Serve()
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
	app := gobastion.New("")
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
type MyFinalizer struct{}

func (f MyFinalizer) Finalize() error {
	log.Printf("[finalizer:MyFinalizer] doing something")
	return nil
}

func main() {
	bastion := gobastion.New("")
	bastion.AppendFinalizers(MyFinalizer{})
	bastion.Serve()
}
```

## Configuration
Represents the configuration for bastion. Config are used to define how the application should run.

###YAML
```yaml
api:
  base_path: "/"
server:
  address: ":8080"

```
###JSON
```json
{
  "api": {
    "base_path": "/"
  },
  "server": {
    "address": ":8080"
  }
}
```

### api
#### `api.base_path`
base path value where the application is going to be mounted. Default `/`.

```json
"base_path": "/foo/test",
```

```
http://localhost/foo/test
```

### `server`
#### `server.address`
Address is the host and port where the app is serve. Default `127.0.0.1:8080`.
When `server.address` is not provided it'll search the ADDR and PORT environment variables 
before set the default.
