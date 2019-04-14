package render

import "net/http"

// StringResponder response strings
type StringResponder string

const (
	// Text returns a StringResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as text/plain.
	Text StringResponder = "text/plain; charset=utf-8"
	// HTML returns a StringResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as text/html.
	HTML StringResponder = "text/html; charset=utf-8"
)

// Response encoded responses in the ResponseWriter with the HTTP status code.
func (s StringResponder) Response(w http.ResponseWriter, code int, response string) {
	writeContentType(w, string(s))
	write(w, code, []byte(response))
}
