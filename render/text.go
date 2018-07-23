package render

import "net/http"

// Text encode the response as "text/plain" content type.
type Text struct {
	w http.ResponseWriter
}

// NewText returns a new Text instance.
func NewText(w http.ResponseWriter) *Text {
	return &Text{w: w}
}

// Response writes a string to the response, setting the Content-Type
// as text/plain.
func (t *Text) Response(code int, v string) {
	t.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	t.w.WriteHeader(code)
	t.w.Write([]byte(v))
}
