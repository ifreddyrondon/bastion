package render

import "net/http"

type StringResponder string

const (
	// Text returns a StringResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as text/plain.
	Text StringResponder = "text/plain; charset=utf-8"
	// HTML returns a StringResponder who's Response writes a string into a ResponseWriter
	// with the Content-Type as text/html.
	HTML StringResponder = "text/html; charset=utf-8"
)

func (s StringResponder) Response(w http.ResponseWriter, code int, response string) {
	w.Header().Set("Content-Type", string(s))
	write(w, code, []byte(response))
}
