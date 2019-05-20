package bastion

import (
	"context"
	"fmt"
	"io"

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
