package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/bastion/render"
)

var internalErrDefaultMsg = errors.New("looks like something went wrong")

// InternalErrCallback sets the callback function when internal error middleware catch a 500 error.
func InternalErrCallback(f func(int, io.Reader)) func(*internalErr) {
	return func(a *internalErr) {
		a.callback = f
	}
}

// InternalErrMsg set default error message to be sent
func InternalErrMsg(err error) func(*internalErr) {
	return func(a *internalErr) {
		a.defaultErr = err
	}
}

type internalErr struct {
	defaultErr error
	render     render.ServerErrRenderer
	callback   func(code int, reader io.Reader)
}

func internalErrCfg(opts ...func(*internalErr)) *internalErr {
	cfg := &internalErr{
		defaultErr: internalErrDefaultMsg,
		render:     render.JSON,
	}
	for _, opt := range opts {
		opt(cfg)
	}
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
					if cfg.callback != nil {
						cfg.callback(m.Code, buf)
					}
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
