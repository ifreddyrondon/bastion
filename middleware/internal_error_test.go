package middleware_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware"

	"gopkg.in/gavv/httpexpect.v1"
)

func TestInternalErrShouldResponseDefaultErrorMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("test"))
	})

	m := middleware.InternalError()
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
}

func TestInternalErrShouldCallTheCallbackFn(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("this should be logged"))
	})

	callback := func(code int, r io.Reader) {
		assert.Equal(t, 500, code)
		var buf bytes.Buffer
		buf.ReadFrom(r)
		assert.Contains(t, buf.String(), "this should be logged")
	}
	m := middleware.InternalError(middleware.InternalErrCallback(callback))
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
}

func TestInternalErrMsg(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("bla"))
	})

	err := errors.New("test")
	m := middleware.InternalError(middleware.InternalErrMsg(err))
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
}

func TestInternalErrNot500(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this should be flushed"))
	})

	callback := func(code int, r io.Reader) {
		assert.Fail(t, "the callback fn should not be called")
	}
	m := middleware.InternalError(middleware.InternalErrCallback(callback))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(200).Body().Equal("this should be flushed")
}
