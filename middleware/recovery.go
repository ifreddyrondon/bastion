package middleware

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/bastion/render"
)

// RecoveryCallback sets the callback function to handler the request when recovers from panics.
func RecoveryCallback(f func(req *http.Request, err error)) func(*recoveryCfg) {
	return func(a *recoveryCfg) {
		a.callback = f
	}
}

type recoveryCfg struct {
	render   render.ServerErrRenderer
	callback func(req *http.Request, err error)
}

func getRecoveryCfg(opts ...func(*recoveryCfg)) *recoveryCfg {
	r := &recoveryCfg{
		render: render.JSON,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(opts ...func(*recoveryCfg)) func(http.Handler) http.Handler {
	cfg := getRecoveryCfg(opts...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch t := r.(type) {
					case error:
						err = errors.WithStack(t)
					case string:
						err = errors.WithStack(errors.New(t))
					default:
						err = errors.New(fmt.Sprint(t))
					}
					if cfg.callback != nil {
						cfg.callback(req, err)
					}
					cfg.render.InternalServerError(w, err)
					return
				}
			}()
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}
