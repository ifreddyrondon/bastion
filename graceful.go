package bastion

import "context"

// graceful shutdown
func graceful(ctx context.Context, app *Bastion) {
	<-ctx.Done()
	logger := app.Logger.With().
		Str("component", "gracefull").
		Logger()
	logger.Info().Msg("preparing for shutdown")
	if err := app.server.Shutdown(ctx); err != nil {
		logger.Error().Err(err)
	} else {
		logger.Info().Msg("gracefully stopped")
	}
}
