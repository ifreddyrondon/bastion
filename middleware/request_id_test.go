package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/middleware"
)

func TestGetRequestIDMissingInstance(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	v := middleware.GetRequestID(ctx)
	assert.Equal(t, "", v)
}

func TestGetRequestIDInvalidReference(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), middleware.RequestIDCtxKey, 1)
	v := middleware.GetRequestID(ctx)
	assert.Equal(t, "", v)
}

func TestGetRequestIDMissingContext(t *testing.T) {
	t.Parallel()
	v := middleware.GetRequestID(nil)
	assert.Equal(t, "", v)
}

func setupReqID(m func(http.Handler) http.Handler) (*httptest.Server, *string, func()) {
	var result string

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetRequestID(r.Context())
		result = reqID
		w.Write([]byte("hi"))
	})

	server := httptest.NewServer(m(h))
	teardown := func() {
		server.Close()
	}
	return server, &result, teardown
}

func TestRequestIDWhenIsNotPresent(t *testing.T) {
	t.Parallel()
	m := middleware.RequestID()
	server, v, teardown := setupReqID(m)
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
	assert.NotNil(t, v)
}

func TestRequestIDWhenIsPresent(t *testing.T) {
	t.Parallel()
	m := middleware.RequestID()
	server, v, teardown := setupReqID(m)
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithHeader("X-Request-Id", "0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusOK)
	assert.Equal(t, "0167c8a5-d308-8692-809d-b1ad4a2d9562", *v)
}

func TestRequestIDWhenIsPresentButIntoAnotherHeader(t *testing.T) {
	t.Parallel()
	m := middleware.RequestID(middleware.RequestIDHeaderName("CHAMO-Id"))
	server, v, teardown := setupReqID(m)
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithHeader("CHAMO-Id", "0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusOK)
	assert.Equal(t, "0167c8a5-d308-8692-809d-b1ad4a2d9562", *v)
}
