package todo_test

import (
	"os"
	"testing"

	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/_examples/todo-rest/todo"
)

var bastion *gobastion.Bastion

func TestMain(m *testing.M) {
	bastion = gobastion.New(nil)
	reader := new(gobastion.JsonReader)
	handler := todo.Handler{
		Reader:    reader,
		Responder: gobastion.DefaultResponder,
	}
	bastion.APIRouter.Mount("/todo/", handler.Routes())
	code := m.Run()
	os.Exit(code)
}

func TestHandlerCreate(t *testing.T) {
	payload := map[string]interface{}{
		"description": "new description",
	}

	e := gobastion.Tester(t, bastion)
	e.POST("/todo/").WithJSON(payload).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 0).
		ContainsKey("description").ValueEqual("description", "new description")
}

func TestHandlerList(t *testing.T) {
	e := gobastion.Tester(t, bastion)
	array := e.GET("/todo/").Expect().
		Status(http.StatusOK).
		JSON().Array().NotEmpty()

	array.Length().Equal(2)
	array.First().Object().
		ContainsKey("id").
		ContainsKey("description")
}

func TestHandlerGet(t *testing.T) {
	e := gobastion.Tester(t, bastion)
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

	e := gobastion.Tester(t, bastion)
	e.PUT("/todo/4").WithJSON(payload).Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 4).
		ContainsKey("description").ValueEqual("description", "updated description")
}

func TestHandlerDelete(t *testing.T) {
	e := gobastion.Tester(t, bastion)
	e.DELETE("/todo/4").Expect().
		Status(http.StatusNoContent).NoContent()
}
