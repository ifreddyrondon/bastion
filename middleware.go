package bastion

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// Recovery is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recovery prints a request ID if one is provided.
func Recovery(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l := logger.With().Str("component", "recovery").Logger()

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
					if err = json.NewRender(w).InternalServerError(err); err != nil {
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
