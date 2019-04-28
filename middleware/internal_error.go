package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/ifreddyrondon/bastion/render"
)

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
		render:       render.JSON,
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
			buf := &bytes.Buffer{}
			writeHeaderHook := WriteHeaderHook(HijackWriteHeaderHook)
			writeHook := WriteHook(HijackWriteHook(buf))
			m, snoop := WrapResponseWriter(w, writeHeaderHook, writeHook)
			defer func() {
				if m.Code >= 500 {
					cfg.logger.Info().
						Str("component", "internal error middleware").
						Int("status", m.Code).
						Msg(buf.String())
					cfg.render.InternalServerError(w, cfg.defaultErr)
					return
				}
				w.WriteHeader(m.Code)
				w.Write(buf.Bytes())
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}
}
