package render

import (
	"net/http"
)

// Engine define methods to encoded response in the body of a request with the HTTP status code.
type Engine interface {
	Response(code int, response interface{}) error
	Send(response interface{}) error
	Created(response interface{}) error
	NoContent()
	BadRequest(err error) error
	NotFound(err error) error
	MethodNotAllowed(err error) error
	InternalServerError(err error) error
}

// Render returns a Engine to response a request with the HTTP status code.
type Render func(http.ResponseWriter) Engine
