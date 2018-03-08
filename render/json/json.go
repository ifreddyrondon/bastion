package json

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
)

// NewRenderer returns the Engine to response with "application/json" content type.
func NewRenderer(w http.ResponseWriter) render.Engine {
	return &Renderer{w}
}

// Renderer encode the response as "application/json" content type and implement the Renderer interface.
type Renderer struct {
	value http.ResponseWriter
}

// Response sends a JSON-encoded response in the body of a request with the HTTP status code.
func (res *Renderer) Response(code int, response interface{}) {
	res.value.Header().Set("Content-Type", "application/json")
	res.value.WriteHeader(code)
	json.NewEncoder(res.value).Encode(response)
}

// Send sends a JSON-encoded response in the body of a request with the 200 status code.
func (res *Renderer) Send(response interface{}) {
	res.Response(http.StatusOK, response)
}

// Created sends a JSON-encoded response in the body of a request with the 201 status code.
func (res *Renderer) Created(response interface{}) {
	res.Response(http.StatusCreated, response)
}

// NoContent sends a response without no content with the 204 status code.
func (res *Renderer) NoContent() {
	res.value.WriteHeader(http.StatusNoContent)
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
	Status  int    `json:"status"`
}

// BadRequest sends a JSON-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (res *Renderer) BadRequest(err error) {
	message := HTTPError{
		Status:  http.StatusBadRequest,
		Errors:  http.StatusText(http.StatusBadRequest),
		Message: err.Error(),
	}
	res.Response(http.StatusBadRequest, message)
}

// NotFound sends a JSON-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (res *Renderer) NotFound(err error) {
	message := HTTPError{
		Status:  http.StatusNotFound,
		Errors:  http.StatusText(http.StatusNotFound),
		Message: err.Error(),
	}
	res.Response(http.StatusNotFound, message)
}

// MethodNotAllowed sends a JSON-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (res *Renderer) MethodNotAllowed(err error) {
	message := HTTPError{
		Status:  http.StatusMethodNotAllowed,
		Errors:  http.StatusText(http.StatusMethodNotAllowed),
		Message: err.Error(),
	}
	res.Response(http.StatusMethodNotAllowed, message)
}

// InternalServerError sends a JSON-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (res *Renderer) InternalServerError(err error) {
	message := HTTPError{
		Status:  http.StatusInternalServerError,
		Errors:  http.StatusText(http.StatusInternalServerError),
		Message: err.Error(),
	}
	res.Response(http.StatusInternalServerError, message)
}