package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

type recoveryCfg struct {
	render       render.ServerErrRenderer
	loggerWriter io.Writer
	logger       zerolog.Logger
}

// RecoveryLoggerOutput set the output for the logger
func RecoveryLoggerOutput(w io.Writer) func(*recoveryCfg) {
	return func(r *recoveryCfg) {
		r.loggerWriter = w
	}
}

func getRecoveryCfg(opts ...func(*recoveryCfg)) *recoveryCfg {
	r := &recoveryCfg{
		render:       render.NewJSON(),
		loggerWriter: os.Stdout,
	}

	for _, opt := range opts {
		opt(r)
	}

	r.logger = zerolog.New(r.loggerWriter).With().Timestamp().Logger()
	return r
}

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(opts ...func(*recoveryCfg)) func(http.Handler) http.Handler {
	confg := getRecoveryCfg(opts...)

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
					confg.logger.Error().
						Str("component", "recovery").
						Err(err).Dict("req", logreq(req)).
						Msg("Recovery middleware catch an error")
					confg.render.InternalServerError(w, err)
					return
				}
			}()
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}
