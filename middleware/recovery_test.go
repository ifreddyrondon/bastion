package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/middleware"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		panicArg interface{}
	}{
		{
			"recovery with err panic call",
			errors.New("testing recovery"),
		},
		{
			"recovery with string panic call",
			"testing recovery",
		},
		{
			"recovery with empty panic call",
			500,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(tc.panicArg)
			})
			m := middleware.Recovery()
			server := httptest.NewServer(m(h))
			defer server.Close()
			expectedRes := map[string]interface{}{
				"message": fmt.Sprintf("%v", tc.panicArg),
				"error":   "Internal Server Error",
				"status":  500,
			}

			e := httpexpect.New(t, server.URL)
			e.GET("/").Expect().Status(500).
				JSON().
				Object().ContainsMap(expectedRes)
		})
	}
}

func TestRecoveryShouldCallTheCallbackFn(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})
	callback := func(req *http.Request, err error) {
		assert.Equal(t, "/", req.URL.RequestURI())
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "HTTP/1.1", req.Proto)
		assert.EqualError(t, err, "test")
	}
	m := middleware.Recovery(middleware.RecoveryCallback(callback))
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
