package todo_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
)

var app *bastion.Bastion

func TestMain(m *testing.M) {
	app = bastion.New(nil)
	reader := new(bastion.JsonReader)
	handler := todo.Handler{
		Reader:    reader,
		Responder: bastion.DefaultResponder,
	}
	app.APIRouter.Mount("/todo/", handler.Routes())
	code := m.Run()
	os.Exit(code)
}

func TestHandlerCreate(t *testing.T) {
	payload := map[string]interface{}{
		"description": "new description",
	}

	e := bastion.Tester(t, app)
	e.POST("/todo/").WithJSON(payload).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 0).
		ContainsKey("description").ValueEqual("description", "new description")
}

func TestHandlerList(t *testing.T) {
	e := bastion.Tester(t, app)
	array := e.GET("/todo/").Expect().
		Status(http.StatusOK).
		JSON().Array().NotEmpty()

	array.Length().Equal(2)
	array.First().Object().
		ContainsKey("id").
		ContainsKey("description")
}

func TestHandlerGet(t *testing.T) {
	e := bastion.Tester(t, app)
	e.GET("/todo/2").Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 2).
		ContainsKey("description").ValueEqual("description", "do something 2")
}

func TestHandlerUpdate(t *testing.T) {
	payload := map[string]interface{}{
		"id":          4,
		"description": "updated description",
	}

	e := bastion.Tester(t, app)
	e.PUT("/todo/4").WithJSON(payload).Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 4).
		ContainsKey("description").ValueEqual("description", "updated description")
}

func TestHandlerDelete(t *testing.T) {
	e := bastion.Tester(t, app)
	e.DELETE("/todo/4").Expect().
		Status(http.StatusNoContent).NoContent()
}
