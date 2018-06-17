package middleware

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/rs/zerolog"
)

// hooks defines WriteHeader and Write methods interceptors for methods included in
// http.ResponseWriter.
func hooks(wc *writerCollector) httpsnoop.Hooks {
	return httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				if !wc.wroteHeader {
					wc.code = code
					wc.wroteHeader = true
				}
			}
		},
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				n, err := wc.buf.Write(p)
				wc.bytes += int64(n)
				wc.wroteHeader = true
				return n, err
			}
		},
	}
}

type writerCollector struct {
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
	// buf store all the []bytes internaly when ResponseWriter.Write is called.
	buf bytes.Buffer
}

func newWriterCollector() *writerCollector {
	return &writerCollector{code: http.StatusOK}
}

// APIErrHandler intercept responses to verify if his status code is >= 500.
// If status is >= 500, it'll response with a default error.
// This middleware allows to response with the same error without disclosure
// internal information, also the real error is logged.
func APIErrHandler(defaultErr error, logger *zerolog.Logger, render render.Render) func(http.Handler) http.Handler {
	l := logger.With().Str("component", "api_error_handler").Logger()

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m := newWriterCollector()

			snoopw := httpsnoop.Wrap(w, hooks(m))
			defer func(logger zerolog.Logger) {
				if m.code >= 500 {
					logger.Error().Int("status", m.code).Bytes("response", m.buf.Bytes()).Msg("")
					if err := render(w).InternalServerError(defaultErr); err != nil {
						logger.Error().Err(err).Msg("")
					}
					return
				}
				w.WriteHeader(m.code)
				if _, err := w.Write(m.buf.Bytes()); err != nil {
					logger.Error().Err(err).Msg("")
				}
			}(l)
			next.ServeHTTP(snoopw, r)
		}
		return http.HandlerFunc(fn)
	}
}
