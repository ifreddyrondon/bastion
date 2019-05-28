package middleware

import (
	"context"
	"net/http"

	"github.com/ifreddyrondon/bastion/ulid"
)

// RequestIDCtxKey is the key that holds the unique request ID in a request context
var RequestIDCtxKey = &contextKey{"Request-id"}

const DefaultRequestIDHeaderName = "X-Request-Id"

// RequestIDHeaderName sets the header name where the request id is searched.
func RequestIDHeaderName(headerName string) func(*requestIDCfg) {
	return func(cfg *requestIDCfg) {
		cfg.headerName = headerName
	}
}

type requestIDCfg struct {
	headerName string
}

func requestIDSetup(opts ...func(*requestIDCfg)) *requestIDCfg {
	cfg := &requestIDCfg{headerName: DefaultRequestIDHeaderName}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// RequestID is a middleware that injects a request ID into the context of each
// request. It search for a request id header or generates one.
// When the request ID is generated, it uses a uuid string that uniquely identifies
// the process.
// The header where the request id is searched can be modified through an option function.
func RequestID(opts ...func(*requestIDCfg)) func(http.Handler) http.Handler {
	cfg := requestIDSetup(opts...)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(cfg.headerName)
			if requestID == "" {
				requestID = ulid.New().String()
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, RequestIDCtxKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// GetRequestID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDCtxKey).(string); ok {
		return reqID
	}
	return ""
}
