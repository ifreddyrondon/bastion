package todo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
)

type Handler struct {
	bastion.Reader
	bastion.Responder
}

// Routes creates a REST router for the todos resource
func (h *Handler) Routes() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", h.List)    // GET /todos - read a list of todos
	r.Post("/", h.Create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", h.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", h.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	todo1 := todo{Description: "do something 1"}
	todo2 := todo{Description: "do something 2"}

	h.Send(w, []todo{todo1, todo2})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var todo1 todo
	if err := h.Reader.Read(r.Body, &todo1); err != nil {
		panic(err) // the error should be handle
	}
	h.Created(w, todo1)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	todo1 := todo{Id: i, Description: fmt.Sprintf("do something %v", id)}
	h.Send(w, todo1)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	var todo1 todo
	if err := h.Reader.Read(r.Body, &todo1); err != nil {
		panic(err)
	}
	todo1.Id = i
	h.Send(w, todo1)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// handle delete logic
	h.NoContent(w)
}
