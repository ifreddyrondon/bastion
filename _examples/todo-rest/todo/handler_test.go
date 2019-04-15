package todo_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/_examples/todo-rest/todo"
)

func setup() *bastion.Bastion {
	app := bastion.New()
	app.Mount("/todo/", todo.Routes())
	return app
}

func TestHandlerCreate(t *testing.T) {
	app := setup()
	payload := map[string]interface{}{"description": "new description"}

	e := bastion.Tester(t, app)
	e.POST("/todo/").WithJSON(payload).Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 0).
		ContainsKey("description").ValueEqual("description", "new description")
}

func TestHandlerList(t *testing.T) {
	app := setup()
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
	app := setup()
	e := bastion.Tester(t, app)
	e.GET("/todo/2").Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", 2).
		ContainsKey("description").ValueEqual("description", "do something 2")
}

func TestHandlerUpdate(t *testing.T) {
	app := setup()
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
	app := setup()
	e := bastion.Tester(t, app)
	e.DELETE("/todo/4").Expect().
		Status(http.StatusNoContent).NoContent()
}

func Test500Err(t *testing.T) {
	app := setup()
	e := bastion.Tester(t, app)

	expectedRes := map[string]interface{}{
		"message": "looks like something went wrong",
		"error":   "Internal Server Error",
		"status":  500,
	}
	e.GET("/todo/error500").Expect().
		Status(http.StatusInternalServerError).
		JSON().
		Object().ContainsMap(expectedRes)
}
