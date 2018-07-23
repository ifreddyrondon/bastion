package render

import "net/http"

// Data encode the response as "application/octet-stream" content type.
type Data struct {
	w http.ResponseWriter
}

// NewData returns a new Data instance.
func NewData(w http.ResponseWriter) *Data {
	return &Data{w: w}
}

// Response writes raw bytes to the response, setting the Content-Type as
// application/octet-stream.
func (t *Data) Response(code int, v []byte) {
	t.w.Header().Set("Content-Type", "application/octet-stream")
	t.w.WriteHeader(code)
	t.w.Write(v)
}
