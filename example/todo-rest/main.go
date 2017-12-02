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
