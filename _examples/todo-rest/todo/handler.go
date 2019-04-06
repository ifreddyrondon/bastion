package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Routes creates a REST router for the todos resource
func Routes() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", list)    // GET /todos - read a list of todos
	r.Post("/", create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	r.Get("/error500", error500) // GET /panic - testing 500 error

	return r
}

func list(w http.ResponseWriter, r *http.Request) {
	todo1 := todo{Description: "do something 1"}
	todo2 := todo{Description: "do something 2"}
	render.NewJSON().Send(w, []todo{todo1, todo2})
}

func create(w http.ResponseWriter, r *http.Request) {
	var todo1 todo
	if err := json.NewDecoder(r.Body).Decode(&todo1); err != nil {
		panic(err) // the error should be handle
	}
	render.NewJSON().Created(w, todo1)
}

func get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	todo1 := todo{ID: i, Description: fmt.Sprintf("do something %v", id)}
	render.NewJSON().Send(w, todo1)
}

func update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	var todo1 todo
	if err := json.NewDecoder(r.Body).Decode(&todo1); err != nil {
		panic(err) // the error should be handle
	}
	todo1.ID = i
	render.NewJSON().Send(w, todo1)
}

func delete(w http.ResponseWriter, r *http.Request) {
	// handle delete logic
	render.NewJSON().NoContent(w)
}

func error500(w http.ResponseWriter, r *http.Request) {
	err := errors.New("test")
	render.NewJSON().InternalServerError(w, err)
}
