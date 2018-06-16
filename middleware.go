package bastion

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// LoggerRequest some provided extra handler to set some request's context fields.
func LoggerRequest(opts *Options) []func(next http.Handler) http.Handler {
	hdls := []func(next http.Handler) http.Handler{}

	access := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		getLoggerWithLevel(r, status).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})

	hdls = append(hdls, access)
	hdls = append(hdls, hlog.RequestIDHandler("req_id", "Request-Id"))
	if !opts.isDEV() {
		hdls = append(hdls, hlog.RemoteAddrHandler("ip"))
		hdls = append(hdls, hlog.UserAgentHandler("user_agent"))
		hdls = append(hdls, hlog.RefererHandler("referer"))
	}
	return hdls
}

func getLoggerWithLevel(r *http.Request, status int) *zerolog.Event {
	if status >= 500 {
		return hlog.FromRequest(r).Error()
	}
	return hlog.FromRequest(r).Info()
}
