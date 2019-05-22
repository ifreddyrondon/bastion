package bastion

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func getLogger(output io.Writer, prettyLogging bool, lvlStr string) (*zerolog.Logger, error) {
	if prettyLogging {
		output = zerolog.ConsoleWriter{Out: output}
	}
	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("bastion logger level unknown: %v", lvlStr))
	}
	l := zerolog.New(output).Level(lvl).With().Timestamp().Logger()
	return &l, nil
}

// LoggerFromCtx returns the Logger associated with the ctx.
func LoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
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
