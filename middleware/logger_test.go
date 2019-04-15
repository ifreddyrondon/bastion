package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/middleware"
)

func TestLoggerDefaults(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	out := &bytes.Buffer{}
	l := zerolog.New(out)
	m := middleware.Logger(middleware.AttachLogger(l))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("ok")
	assert.Contains(t, out.String(), `"status":200`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"URL":"/"`)
	assert.Contains(t, out.String(), `"size"`)
	assert.Contains(t, out.String(), `"duration"`)
	assert.Contains(t, out.String(), `"level":"info`)
	assert.Contains(t, out.String(), `"req_id"`)
	assert.NotContains(t, out.String(), `"ip"`)
	assert.NotContains(t, out.String(), `"user_agent"`)
	assert.NotContains(t, out.String(), `"referer"`)
}

func TestLoggerRequestLevelErrorForStatusGreaterThan500(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}
	l := zerolog.New(out)
	m := middleware.Logger(middleware.AttachLogger(l))
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
	server := httptest.NewServer(m(h))
	defer server.Close()
	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(500)
	assert.Contains(t, out.String(), `"status":500`)
	assert.Contains(t, out.String(), `"method":"GET"`)
	assert.Contains(t, out.String(), `"level":"error`)
}

func TestLoggerOptions(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	out := &bytes.Buffer{}
	l := zerolog.New(out)
	m := middleware.Logger(
		middleware.AttachLogger(l),
		middleware.DisableLogDuration(),
		middleware.DisableLogMethod(),
		middleware.DisableLogRequestID(),
		middleware.DisableLogSize(),
		middleware.DisableLogStatus(),
		middleware.DisableLogURL(),
		middleware.EnableLogReferer(),
		middleware.EnableLogReqIP(),
		middleware.EnableLogUserAgent(),
	)
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("ok")
	assert.Contains(t, out.String(), `"level":"info`)
	assert.NotContains(t, out.String(), `"status":200`)
	assert.NotContains(t, out.String(), `"method":"GET"`)
	assert.NotContains(t, out.String(), `"URL":"/"`)
	assert.NotContains(t, out.String(), `"size"`)
	assert.NotContains(t, out.String(), `"duration"`)
	assert.NotContains(t, out.String(), `"req_id"`)
	assert.Contains(t, out.String(), `"ip"`)
	assert.Contains(t, out.String(), `"user_agent"`)
}

func TestLoggerRequestErrorLvl(t *testing.T) {
	t.Parallel()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	out := &bytes.Buffer{}
	l := zerolog.New(out).Level(zerolog.Level(zerolog.ErrorLevel))
	m := middleware.Logger(middleware.AttachLogger(l))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("ok")
	assert.NotContains(t, out.String(), `"status":200`)
}
