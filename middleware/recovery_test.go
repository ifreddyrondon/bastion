package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
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

			out := &bytes.Buffer{}
			m := middleware.Recovery(middleware.RecoveryLoggerOutput(out))
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

			output := out.String()
			assert.Contains(t, output, `"component":"recovery`)
			assert.Contains(t, output, `"message":"Recovery middleware catch an error`)
		})
	}
}

func TestRecoveryLogRequestGET(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	m := middleware.Recovery(middleware.RecoveryLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"GET","proto":"HTTP/1.1","host":"`)
	assert.Contains(t, out.String(), `"headers":{"user-agent":"Go-http-client/1.1","accept-encoding":"gzip"},"body":""}`)

}

func TestRecoveryLogRequestWithHeaders(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	m := middleware.Recovery(middleware.RecoveryLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	e.GET("/").WithHeader("User-Agent", "Mozilla").Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"GET","proto":"HTTP/1.1","host":"`)
	assert.Contains(t, out.String(), `"headers":{"user-agent":"Mozilla","accept-encoding":"gzip"},"body":""}`)
}

func TestRecoveryLogRequestPOST(t *testing.T) {
	t.Parallel()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	out := &bytes.Buffer{}
	m := middleware.Recovery(middleware.RecoveryLoggerOutput(out))
	server := httptest.NewServer(m(h))
	defer server.Close()

	payload := map[string]string{"hello": "world"}
	e := httpexpect.New(t, server.URL)
	e.POST("/").WithJSON(payload).
		Expect().Status(500).JSON()
	assert.Contains(t, out.String(), `"level":"error`)
	assert.Contains(t, out.String(), `"component":"recovery"`)
	assert.Contains(t, out.String(), `"error":"test"`)
	assert.Contains(t, out.String(), `"req":{"url":"/","method":"POST","proto":"HTTP/1.1","host":"`)
	assert.Contains(t, out.String(), `"body":"{\"hello\":\"world\"}"`)
}
