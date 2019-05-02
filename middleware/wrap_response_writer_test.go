package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/middleware"
)

func TestWrapResponseWriterDefaultHooks(t *testing.T) {
	t.Parallel()

	var metrics middleware.WriterMetricsCollector
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("created"))
	})

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m, snoop := middleware.WrapResponseWriter(w)
			defer func() {
				metrics = *m
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}

	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(201).Body().Equal("created")

	assert.Equal(t, 201, metrics.Code)
	assert.Equal(t, int64(7), metrics.Bytes)
}

func TestWrapResponseWriterDefaultWithOverwriteWriteHeaderHook(t *testing.T) {
	t.Parallel()

	var metrics middleware.WriterMetricsCollector
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("created"))
	})

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m, snoop := middleware.WrapResponseWriter(
				w,
				middleware.WriteHeaderHook(middleware.HijackWriteHeaderHook),
			)
			defer func() {
				metrics = *m
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}

	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("created")

	assert.Equal(t, 201, metrics.Code)
	assert.Equal(t, int64(7), metrics.Bytes)
}

func TestWrapResponseWriterDefaultWithCopyWriterHook(t *testing.T) {
	t.Parallel()

	var metrics middleware.WriterMetricsCollector
	var out bytes.Buffer
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("created"))
	})

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m, snoop := middleware.WrapResponseWriter(
				w,
				middleware.WriteHook(middleware.CopyWriterHook(&out)),
			)
			defer func() {
				metrics = *m
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}

	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(201).Body().Equal("created")

	assert.Equal(t, 201, metrics.Code)
	assert.Equal(t, int64(7), metrics.Bytes)
	assert.Equal(t, out.String(), "created")
}

func TestWrapResponseWriterDefaultWithOverwriteWriteHook(t *testing.T) {
	t.Parallel()

	var metrics middleware.WriterMetricsCollector
	var out bytes.Buffer
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("created"))
	})

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m, snoop := middleware.WrapResponseWriter(
				w,
				middleware.WriteHook(middleware.HijackWriteHook(&out)),
			)
			defer func() {
				metrics = *m
			}()
			next.ServeHTTP(snoop, r)
		}
		return http.HandlerFunc(fn)
	}

	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(201).Body().Equal("")

	assert.Equal(t, 201, metrics.Code)
	assert.Equal(t, int64(7), metrics.Bytes)
	assert.Equal(t, out.String(), "created")
}
