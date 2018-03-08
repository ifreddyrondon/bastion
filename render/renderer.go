package render

import (
	"net/http"
)

// RendererEngine define methods to encoded response in the body of a request with the HTTP status code.
type RendererEngine interface {
	Response(code int, response interface{})
	Send(response interface{})
	Created(response interface{})
	NoContent()
	BadRequest(err error)
	NotFound(err error)
	MethodNotAllowed(err error)
	InternalServerError(err error)
}

// Renderer returns a RendererEngine to response a request with the HTTP status code.
type Renderer func(http.ResponseWriter) RendererEngine
