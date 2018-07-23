package render

import "net/http"

// HTML encode the response as "text/html; charset=utf-8" content type.
type HTML struct {
	w http.ResponseWriter
}

// NewHTML returns a new HTML instance.
func NewHTML(w http.ResponseWriter) *HTML {
	return &HTML{w: w}
}

// Response writes a string to the response, setting the Content-Type as text/html.
func (t *HTML) Response(code int, v string) {
	t.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.w.WriteHeader(code)
	t.w.Write([]byte(v))
}
