package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/felixge/httpsnoop"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/pkg/errors"
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

// ErrAPIDefault default error message response when something happens.
var ErrAPIDefault = errors.New("looks like something went wrong")

type apiErrCfg struct {
	defaultErr   error
	render       render.ServerErrRenderer
	loggerWriter io.Writer
	logger       zerolog.Logger
}

// APIErrorLoggerOutput set the output for the logger
func APIErrorLoggerOutput(w io.Writer) func(*apiErrCfg) {
	return func(a *apiErrCfg) {
		a.loggerWriter = w
	}
}

// APIErrorDefault500 set default error message to be sent
func APIErrorDefault500(err error) func(*apiErrCfg) {
	return func(a *apiErrCfg) {
		a.defaultErr = err
	}
}

func getAPIErrCfg(opts ...func(*apiErrCfg)) *apiErrCfg {
	a := &apiErrCfg{
		defaultErr:   ErrAPIDefault,
		render:       render.NewJSON(),
		loggerWriter: os.Stdout,
	}

	for _, opt := range opts {
		opt(a)
	}

	a.logger = zerolog.New(a.loggerWriter).With().Timestamp().Logger()
	return a
}

// APIError intercept responses to verify their status and handle the error.
// It gets the response code and if it's >= 500 handlers the error with a
// default error message and without disclosure internal information.
// The real error keeds logged.
func APIError(opts ...func(*apiErrCfg)) func(http.Handler) http.Handler {
	confg := getAPIErrCfg(opts...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m := newWriterCollector()

			snoopw := httpsnoop.Wrap(w, hooks(m))
			defer func() {
				if m.code >= 500 {
					confg.logger.Info().
						Str("component", "api_error_handler").
						Bytes("response", m.buf.Bytes()).
						Msg("APIError middleware catch a response error >= 500")
					confg.render.InternalServerError(w, confg.defaultErr)
					return
				}
				w.WriteHeader(m.code)
				w.Write(m.buf.Bytes())
			}()
			next.ServeHTTP(snoopw, r)
		}
		return http.HandlerFunc(fn)
	}
}
