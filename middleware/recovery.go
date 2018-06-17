package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func logreq(r *http.Request) *zerolog.Event {
	evt := zerolog.Dict()
	evt.Str("url", r.URL.RequestURI()).
		Str("method", r.Method).
		Str("proto", r.Proto).
		Str("host", r.Host)

	headers := zerolog.Dict()
	for name, values := range r.Header {
		name = strings.ToLower(name)
		headers.Str(name, strings.Join(values, ","))
	}
	evt.Dict("headers", headers)

	if r.Body != nil {
		body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		evt.Bytes("body", body)
	}

	return evt
}

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(logger *zerolog.Logger, render render.Render) func(http.Handler) http.Handler {
	l := logger.With().Str("component", "recovery").Logger()

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			defer func(logger zerolog.Logger) {
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

					logger.Error().Err(err).Dict("req", logreq(req)).Msg("")
					if err = render(w).InternalServerError(err); err != nil {
						logger.Error().Err(err).Msg("")
					}
					return
				}
			}(l)
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}
