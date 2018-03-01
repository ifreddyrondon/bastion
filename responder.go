package bastion

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	Response(w http.ResponseWriter, code int, response interface{})
	Send(w http.ResponseWriter, response interface{})
	Created(w http.ResponseWriter, response interface{})
	NoContent(w http.ResponseWriter)
	BadRequest(w http.ResponseWriter, err error)
	NotFound(w http.ResponseWriter, err error)
	MethodNotAllowed(w http.ResponseWriter, err error)
	InternalServerError(w http.ResponseWriter, err error)
}

// JsonResponder response with the error context.
type JsonResponder struct{}

// Response, sends a JSON-encoded response in the body of a request with the HTTP status code.
func (res *JsonResponder) Response(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ResponseJson, sends a JSON-encoded response in the body of a request with the 200 status code.
func (res *JsonResponder) Send(w http.ResponseWriter, response interface{}) {
	res.Response(w, http.StatusOK, response)
}

// Created, sends a JSON-encoded response in the body of a request with the 201 status code.
func (res *JsonResponder) Created(w http.ResponseWriter, response interface{}) {
	res.Response(w, http.StatusCreated, response)
}

// NoContent, sends a response without no content with the 204 status code.
func (res *JsonResponder) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
	Status  int    `json:"status"`
}

// BadRequest, sends a JSON-encoded error response in the body of a request with the 400 status code.
// The response will contains the status 400 and error "Bad Request".
func (res *JsonResponder) BadRequest(w http.ResponseWriter, err error) {
	message := HTTPError{
		Status:  http.StatusBadRequest,
		Errors:  http.StatusText(http.StatusBadRequest),
		Message: err.Error(),
	}
	res.Response(w, http.StatusBadRequest, message)
}

// Abort, sends a JSON-encoded error response in the body of a request with the 404 status code.
// The response will contains the status 404 and error "Not Found".
func (res *JsonResponder) NotFound(w http.ResponseWriter, err error) {
	message := HTTPError{
		Status:  http.StatusNotFound,
		Errors:  http.StatusText(http.StatusNotFound),
		Message: err.Error(),
	}
	res.Response(w, http.StatusNotFound, message)
}

// MethodNotAllowed, sends a JSON-encoded error response in the body of a request with the 405 status code.
// The response will contains the status 405 and error "Method Not Allowed".
func (res *JsonResponder) MethodNotAllowed(w http.ResponseWriter, err error) {
	message := HTTPError{
		Status:  http.StatusMethodNotAllowed,
		Errors:  http.StatusText(http.StatusMethodNotAllowed),
		Message: err.Error(),
	}
	res.Response(w, http.StatusMethodNotAllowed, message)
}

// InternalServerError, sends a JSON-encoded error response in the body of a request with the 500 status code.
// The response will contains the status 500 and error "Internal Server Error".
func (res *JsonResponder) InternalServerError(w http.ResponseWriter, err error) {
	message := HTTPError{
		Status:  http.StatusInternalServerError,
		Errors:  http.StatusText(http.StatusInternalServerError),
		Message: err.Error(),
	}
	res.Response(w, http.StatusInternalServerError, message)
}
