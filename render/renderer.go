package render

// Renderer interface for managing general response payloads.
type Renderer interface {
	// Response encoded responses in the body of a request with the HTTP status code.
	Response(code int, response interface{})
}

// APIRenderer interface for managing API response payloads.
type APIRenderer interface {
	Renderer
	Send(response interface{})
	Created(response interface{})
	NoContent()
	BadRequest(err error)
	NotFound(err error)
	MethodNotAllowed(err error)
	InternalServerError(err error)
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// NewHTTPError returns a new HTTPError instance.
func NewHTTPError(message, err string, status int) *HTTPError {
	return &HTTPError{
		Message: message,
		Error:   err,
		Status:  status,
	}
}
