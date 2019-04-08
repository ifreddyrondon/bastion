package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// AttachLogger chain the logger with the middleware.
func AttachLogger(log zerolog.Logger) func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.logger = log
	}
}

// EnableLogReqIP show the request ip.
func EnableLogReqIP() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.enableLogReqIP = true
	}
}

// EnableLogUserAgent show the user agent of the request.
func EnableLogUserAgent() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.enableLogUserAgent = true
	}
}

// EnableLogReferer show referer of the request.
func EnableLogReferer() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.enableLogReferer = true
	}
}

// DisableLogMethod hide the request method.
func DisableLogMethod() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogMethod = true
	}
}

// DisableLogURL hide the request url.
func DisableLogURL() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogURL = true
	}
}

// DisableLogStatus hide the request status.
func DisableLogStatus() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogStatus = true
	}
}

// DisableLogStatus hide the request size.
func DisableLogSize() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogSize = true
	}
}

// DisableLogStatus hide the request duration.
func DisableLogDuration() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogDuration = true
	}
}

// DisableLogStatus hide the request id.
func DisableLogRequestID() func(*loggerCfg) {
	return func(r *loggerCfg) {
		r.disableLogRequestID = true
	}
}

type loggerCfg struct {
	logger              zerolog.Logger
	disableLogMethod    bool
	disableLogURL       bool
	disableLogStatus    bool
	disableLogSize      bool
	disableLogDuration  bool
	disableLogRequestID bool
	enableLogReqIP      bool
	enableLogUserAgent  bool
	enableLogReferer    bool
}

func getLoggerCfg(opts ...func(*loggerCfg)) *loggerCfg {
	r := &loggerCfg{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func getLoggerWithLevel(r *http.Request, status int) *zerolog.Event {
	if status >= 500 {
		return hlog.FromRequest(r).Error()
	}
	return hlog.FromRequest(r).Info()
}

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
//
// Alternatively, look at https://github.com/rs/zerolog#integration-with-nethttp.
func Logger(opts ...func(*loggerCfg)) func(http.Handler) http.Handler {
	cfg := getLoggerCfg(opts...)
	loggers := []func(http.Handler) http.Handler{
		hlog.NewHandler(cfg.logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			l := getLoggerWithLevel(r, status)
			if !cfg.disableLogMethod {
				l.Str("method", r.Method)
			}
			if !cfg.disableLogURL {
				l.Str("URL", r.URL.String())
			}
			if !cfg.disableLogStatus {
				l.Int("status", status)
			}
			if !cfg.disableLogSize {
				l.Int("size", size)
			}
			if !cfg.disableLogDuration {
				l.Dur("duration", duration)
			}
			l.Msg("")
		}),
	}
	if !cfg.disableLogRequestID {
		loggers = append(loggers, hlog.RequestIDHandler("req_id", "Request-Id"))
	}
	if cfg.enableLogReqIP {
		loggers = append(loggers, hlog.RemoteAddrHandler("ip"))
	}
	if cfg.enableLogUserAgent {
		loggers = append(loggers, hlog.UserAgentHandler("user_agent"))
	}
	if cfg.enableLogReferer {
		loggers = append(loggers, hlog.RefererHandler("referer"))
	}

	return func(next http.Handler) http.Handler {
		h := loggers[len(loggers)-1](next)
		for i := len(loggers) - 2; i >= 0; i-- {
			h = loggers[i](h)
		}
		return h
	}
}
