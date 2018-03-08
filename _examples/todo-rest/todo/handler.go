package todo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/reader"
	"github.com/ifreddyrondon/bastion/renderer"
)

type Handler struct {
	Reader reader.Reader
	Render renderer.Renderer
}

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
	todo1 := todo{Description: "do something 1"}
	todo2 := todo{Description: "do something 2"}
	h.Render(w).Send([]todo{todo1, todo2})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var todo1 todo
	if err := h.Reader(r.Body).Read(&todo1); err != nil {
		panic(err) // the error should be handle
	}
	h.Render(w).Created(todo1)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	todo1 := todo{Id: i, Description: fmt.Sprintf("do something %v", id)}
	h.Render(w).Send(todo1)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	var todo1 todo
	if err := h.Reader(r.Body).Read(&todo1); err != nil {
		panic(err)
	}
	todo1.Id = i
	h.Render(w).Send(todo1)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	// handle delete logic
	h.Render(w).NoContent()
}
