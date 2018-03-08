package renderer

import (
	"net/http"
)

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
type Renderer func(http.ResponseWriter) Engine
