package middleware_test

import "github.com/pkg/errors"

type mockRenderEngine struct{}

func (m *mockRenderEngine) Response(code int, response interface{}) error { return nil }
func (m *mockRenderEngine) Send(response interface{}) error               { return nil }
func (m *mockRenderEngine) Created(response interface{}) error            { return nil }
func (m *mockRenderEngine) NoContent()                                    {}
func (m *mockRenderEngine) BadRequest(err error) error                    { return nil }
func (m *mockRenderEngine) NotFound(err error) error                      { return nil }
func (m *mockRenderEngine) MethodNotAllowed(err error) error              { return nil }
func (m *mockRenderEngine) InternalServerError(err error) error {
	return errors.New("error render")
}
