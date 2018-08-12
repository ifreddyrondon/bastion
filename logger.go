package bastion

import (
	"context"

	"github.com/rs/zerolog"
)

func getLogger(opts *Options) *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: opts.LoggerOutput}).With().
		Timestamp().
		Str("app", "bastion").
		Logger()

	logger = logger.Level(zerolog.Level(opts.LoggerLevel))
	if opts.NoPrettyLogging {
		logger = logger.Output(opts.LoggerOutput)
	}

	return &logger
}

// LoggerFromCtx returns the Logger associated with the ctx.
func LoggerFromCtx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
