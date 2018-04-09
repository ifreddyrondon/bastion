package bastion

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func getLogger(opts *Options) *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: opts.LoggerWriter}).With().
		Timestamp().
		Str("app", "bastion").
		Logger()

	logger = logger.Level(zerolog.Level(opts.LoggerLevel))
	if opts.NoPrettyLogging {
		logger = logger.Output(opts.LoggerWriter)
	}

	return &logger
}

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

// LoggerFromCtx returns the Logger associated with the ctx.
func LoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
