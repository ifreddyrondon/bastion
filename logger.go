package bastion

import (
	"context"

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

// LoggerFromCtx returns the Logger associated with the ctx.
func LoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
