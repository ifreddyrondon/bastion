package middleware_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware"

	"gopkg.in/gavv/httpexpect.v1"
)

func TestAPIErrCatch500DefaultMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	m := middleware.APIError(middleware.APIErrorLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()
	expectedRes := map[string]interface{}{
		"message": "looks like something went wrong",
		"error":   "Internal Server Error",
		"status":  500,
	}

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(500).
		JSON().
		Object().ContainsMap(expectedRes)

	output := out.String()
	assert.Contains(t, output, `"component":"api_error_handler`)
	assert.Contains(t, output, `"response":"this should be logged`)
	assert.Contains(t, output, `"message":"APIError middleware catch a response error >= 500`)
}

func TestAPIErrCatch500CustomMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	err := errors.New("test")
	m := middleware.APIError(middleware.APIErrorLoggerOutput(out), middleware.APIErrorDefault500(err))
	server := httptest.NewServer(m(h))
	defer server.Close()
	expectedRes := map[string]interface{}{
		"message": "test",
		"error":   "Internal Server Error",
		"status":  500,
	}

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(500).
		JSON().
		Object().ContainsMap(expectedRes)

	output := out.String()
	assert.Contains(t, output, `"component":"api_error_handler`)
	assert.Contains(t, output, `"response":"this should be logged`)
	assert.Contains(t, output, `"message":"APIError middleware catch a response error >= 500`)
}

func TestAPIErrNot500(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this should be flushed"))
	})

	out := &bytes.Buffer{}
	m := middleware.APIError(middleware.APIErrorLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("this should be flushed")

	assert.NotContains(t, out.String(), `"component":"api_error_handler`)
}
