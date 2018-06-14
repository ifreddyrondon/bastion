package middleware

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/rs/zerolog"
)

// Metrics holds metrics captured from CaptureMetrics.
type basicWriter struct {
	// Code is the first http response code passed to the WriteHeader func of
	// the ResponseWriter. If no such call is made, a default code of 200 is
	// assumed instead.
	code int
	// Duration is the time it took to execute the handler.
	// Duration time.Duration
	wroteHeader bool
	// bytes is the number of bytes successfully written by the Write or
	// ReadFrom function of the ResponseWriter. ResponseWriters may also write
	// data to their underlaying connection directly (e.g. headers), but those
	// are not tracked. Therefor the number of Written bytes will usually match
	// the size of the response body.
	bytes int64

	buf bytes.Buffer

	a bool
}

func (b *basicWriter) Reset() {
	b.buf.Reset()
	b.bytes = 0
}

func (b *basicWriter) Unlock() {
	b.a = true
}

// APIErrHandler is the container for error handler middleware
type APIErrHandler struct {
	defaultErr error
	logger     *zerolog.Logger
}

// NewAPIErrHandler returns a new instance of APIErrorHandler
func NewAPIErrHandler(defaultErr error, logger *zerolog.Logger) *APIErrHandler {
	return &APIErrHandler{defaultErr: defaultErr, logger: logger}
}

// Handler intercept responses to verify if his status code is >= 500.
// If status is >= 500, it'll response with a default error.
// This middleware allows to response with the same error without disclosure
// internal information, also the real error is logged.
func (a *APIErrHandler) Handler(next http.Handler) http.Handler {
	l := a.logger.With().Str("component", "api_error_handler").Logger()
	fn := func(w http.ResponseWriter, r *http.Request) {

		m := basicWriter{code: http.StatusOK}
		hooks := httpsnoop.Hooks{
			WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					next(code)
					if !m.wroteHeader {
						m.code = code
						m.wroteHeader = true
					}
				}
			},
			Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(p []byte) (int, error) {
					n, err := m.buf.Write(p)
					m.bytes += int64(n)
					if m.a {
						_, err2 := next(m.buf.Bytes())
						// Prefer errors generated by the proxied writer.
						if err2 == nil {
							err = err2
						}
					}
					m.wroteHeader = true
					return n, err
				}
			},
			// Flush: func(next httpsnoop.FlushFunc) httpsnoop.FlushFunc {

			// },
		}

		snoopw := httpsnoop.Wrap(w, hooks)
		defer func(logger zerolog.Logger) {
			// if m.code >= 500 {
			logger.Error().Int("status", m.code).Bytes("response", m.buf.Bytes()).Msg("")
			m.Reset()
			m.Unlock()
			if err := json.NewRender(snoopw).InternalServerError(a.defaultErr); err != nil {
				logger.Error().Err(err).Msg("")
			}
			return
			// }
			m.Unlock()
			snoopw.Write([]byte(""))
		}(l)
		// next.ServeHTTP(lw, r)
		next.ServeHTTP(snoopw, r)
	}
	return http.HandlerFunc(fn)
}
