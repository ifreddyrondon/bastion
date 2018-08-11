package bastion

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// ProductionLoggers is a list of logger added for production
var ProductionLoggers = []func(next http.Handler) http.Handler{
	hlog.RemoteAddrHandler("ip"),
	hlog.UserAgentHandler("user_agent"),
	hlog.RefererHandler("referer"),
}

func getLoggerWithLevel(r *http.Request, status int) *zerolog.Event {
	if status >= 500 {
		return hlog.FromRequest(r).Error()
	}
	return hlog.FromRequest(r).Info()
}

// Logger middleware to log request.
func loggerRequest(isProd bool) []func(next http.Handler) http.Handler {
	access := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		getLoggerWithLevel(r, status).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})
	hdls := []func(next http.Handler) http.Handler{
		access,
		hlog.RequestIDHandler("req_id", "Request-Id"),
	}
	if isProd {
		hdls = append(hdls, ProductionLoggers...)
	}
	return hdls
}
