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

func TestInternalErrDefaultMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	m := middleware.InternalError(middleware.InternalErrLoggerOutput(out))
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
	assert.Contains(t, output, `"component":"internal error middleware`)
	assert.Contains(t, output, `"status":500`)
	assert.Contains(t, output, `"message":"this should be logged`)
}

func TestInternalErrMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	out := &bytes.Buffer{}
	err := errors.New("test")
	m := middleware.InternalError(middleware.InternalErrLoggerOutput(out), middleware.InternalErrMsg(err))
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
	assert.Contains(t, output, `"component":"internal error middleware`)
	assert.Contains(t, output, `"status":500`)
	assert.Contains(t, output, `"message":"this should be logged`)
}

func TestInternalErrNot500(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this should be flushed"))
	})

	out := &bytes.Buffer{}
	m := middleware.InternalError(middleware.InternalErrLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("this should be flushed")

	assert.NotContains(t, out.String(), `"component":"internal error middleware`)
}
