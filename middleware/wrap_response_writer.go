package middleware

import (
	"io"
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
)

var (
	// CollectHeaderHook capture the response code into a WriterMetricsCollector and forward the execution to the main
	// WriteHeader method. It's the default hook for WriteHeader when WrapWriter is used.
	CollectHeaderHook = func(collector *WriterMetricsCollector) func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
		return func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				next(code)
				collector.locker.Lock()
				defer collector.locker.Unlock()
				if !collector.wroteHeader {
					collector.Code = code
					collector.wroteHeader = true
				}
			}
		}
	}

	// HijackWriteHeaderHook capture the response code into a WriterMetricsCollector. Warning it'll not forward to the
	// main WriteHeader method execution.
	HijackWriteHeaderHook = func(collector *WriterMetricsCollector) func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
		return func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				collector.locker.Lock()
				defer collector.locker.Unlock()
				if !collector.wroteHeader {
					collector.Code = code
					collector.wroteHeader = true
				}
			}
		}
	}

	// CollectBytesHook capture the amount of bytes into a WriterMetricsCollector and forward the execution to the main
	// Write method. It's the default hook for Write when WrapWriter is used.
	CollectBytesHook = func(collector *WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
		return func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				n, err := next(p)
				collector.locker.Lock()
				defer collector.locker.Unlock()
				collector.Bytes += int64(n)
				collector.wroteHeader = true
				return n, err
			}
		}
	}
)

// CopyWriterHook makes a copy of the bytes into a io.Writer and forward the execution to the main Write method.
// It'll save the amount of bytes into a WriterMetricsCollector.
func CopyWriterHook(w io.Writer) func(collector *WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
	return func(collector *WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
		return func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				n, err := next(p)
				collector.locker.Lock()
				defer collector.locker.Unlock()
				w.Write(p)
				collector.Bytes += int64(n)
				collector.wroteHeader = true
				return n, err
			}
		}
	}
}

// HijackWriteHook write the response bytes into a io.Writer. It'll save the amount of bytes into
// a WriterMetricsCollector. Warning it'll not forward to the main Write method execution.
func HijackWriteHook(w io.Writer) func(collector *WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
	return func(collector *WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
		return func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				collector.locker.Lock()
				defer collector.locker.Unlock()
				n, err := w.Write(p)
				collector.Bytes += int64(n)
				collector.wroteHeader = true
				return n, err
			}
		}
	}
}

// WriterMetricsCollector holds metrics captured from writer.
type WriterMetricsCollector struct {
	// Code is the first http response code passed to the WriteHeader func of
	// the ResponseWriter. If no such call is made, a default code of 200 is
	// assumed instead.
	Code int
	// bytes is the number of bytes successfully written by the Write or
	// ReadFrom function of the ResponseWriter. ResponseWriters may also write
	// data to their underlying connection directly (e.g. headers), but those
	// are not tracked. Therefor the number of Written bytes will usually match
	// the size of the response body.
	Bytes       int64
	wroteHeader bool
	locker      sync.Mutex
}

// WriteHeaderHook define the method interceptor when WriteHeader is called.
func WriteHeaderHook(hook func(*WriterMetricsCollector) func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc) func(*wrapWriterOpts) {
	return func(opts *wrapWriterOpts) {
		opts.writeHeaderHook = hook
	}
}

// WriteHook define the method interceptor when Write is called.
func WriteHook(hook func(*WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc) func(*wrapWriterOpts) {
	return func(opts *wrapWriterOpts) {
		opts.writeHook = hook
	}
}

type wrapWriterOpts struct {
	writeHeaderHook func(*WriterMetricsCollector) func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc
	writeHook       func(*WriterMetricsCollector) func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc
}

func wrapWriterSetupCfg(opts ...func(*wrapWriterOpts)) *wrapWriterOpts {
	cfg := &wrapWriterOpts{
		writeHeaderHook: CollectHeaderHook,
		writeHook:       CollectBytesHook,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WrapResponseWriter defines a set of method interceptors for methods included in
// http.ResponseWriter as well as some others. You can think of them as
// middleware for the function calls they target.
// It response with a wrapped ResponseWriter and WriterMetricsCollector with some useful metrics.
func WrapResponseWriter(w http.ResponseWriter, opts ...func(*wrapWriterOpts)) (*WriterMetricsCollector, http.ResponseWriter) {
	cfg := wrapWriterSetupCfg(opts...)
	collector := &WriterMetricsCollector{Code: http.StatusOK}
	hooks := httpsnoop.Hooks{
		WriteHeader: cfg.writeHeaderHook(collector),
		Write:       cfg.writeHook(collector),
	}
	snoop := httpsnoop.Wrap(w, hooks)
	return collector, snoop
}
