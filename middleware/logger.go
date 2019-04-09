package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// AttachLogger chain the logger with the middleware.
func AttachLogger(log zerolog.Logger) LoggerOpt {
	return func(r *loggerCfg) {
		r.logger = log
	}
}

// EnableLogReqIP show the request ip.
func EnableLogReqIP() LoggerOpt {
	return func(r *loggerCfg) {
		r.enableLogReqIP = true
	}
}

// EnableLogUserAgent show the user agent of the request.
func EnableLogUserAgent() LoggerOpt {
	return func(r *loggerCfg) {
		r.enableLogUserAgent = true
	}
}

// EnableLogReferer show referer of the request.
func EnableLogReferer() LoggerOpt {
	return func(r *loggerCfg) {
		r.enableLogReferer = true
	}
}

// DisableLogMethod hide the request method.
func DisableLogMethod() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogMethod = true
	}
}

// DisableLogURL hide the request url.
func DisableLogURL() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogURL = true
	}
}

// DisableLogStatus hide the request status.
func DisableLogStatus() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogStatus = true
	}
}

// DisableLogStatus hide the request size.
func DisableLogSize() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogSize = true
	}
}

// DisableLogStatus hide the request duration.
func DisableLogDuration() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogDuration = true
	}
}

// DisableLogStatus hide the request id.
func DisableLogRequestID() LoggerOpt {
	return func(r *loggerCfg) {
		r.disableLogRequestID = true
	}
}

type LoggerOpt func(*loggerCfg)

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

func getLoggerCfg(opts ...LoggerOpt) *loggerCfg {
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
func Logger(opts ...LoggerOpt) func(http.Handler) http.Handler {
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
