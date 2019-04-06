package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/felixge/httpsnoop"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/ifreddyrondon/bastion/render"
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
	// data to their underlying connection directly (e.g. headers), but those
	// are not tracked. Therefor the number of Written bytes will usually match
	// the size of the response body.
	bytes int64
	// buf store all the []bytes internally when ResponseWriter.Write is called.
	buf bytes.Buffer
}

func newWriterCollector() *writerCollector {
	return &writerCollector{code: http.StatusOK}
}

var internalErrDefaultMsg = errors.New("looks like something went wrong")

// InternalErrLoggerOutput set the output for the logger
func InternalErrLoggerOutput(w io.Writer) func(*internalErr) {
	return func(a *internalErr) {
		a.loggerWriter = w
	}
}

// InternalErrMsg set default error message to be sent
func InternalErrMsg(err error) func(*internalErr) {
	return func(a *internalErr) {
		a.defaultErr = err
	}
}

type internalErr struct {
	defaultErr   error
	render       render.ServerErrRenderer
	loggerWriter io.Writer
	logger       zerolog.Logger
}

func internalErrCfg(opts ...func(*internalErr)) *internalErr {
	cfg := &internalErr{
		defaultErr:   internalErrDefaultMsg,
		render:       render.NewJSON(),
		loggerWriter: os.Stdout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cfg.logger = zerolog.New(cfg.loggerWriter).With().Timestamp().Logger()
	return cfg
}

// InternalError intercept responses to verify their status and handle the error.
// It gets the response code and if it's >= 500 handles the error with a
// default error message without disclosure internal information.
// The real error keeps logged.
func InternalError(opts ...func(*internalErr)) func(http.Handler) http.Handler {
	cfg := internalErrCfg(opts...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m := newWriterCollector()

			snoop := httpsnoop.Wrap(w, hooks(m))
			defer func() {
				if m.code >= 500 {
					cfg.logger.Info().
						Str("component", "internal error middleware").
						Int("status", m.code).
						Msg(m.buf.String())
					cfg.render.InternalServerError(w, cfg.defaultErr)
					return
				}
				w.WriteHeader(m.code)
				w.Write(m.buf.Bytes())
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}
}
