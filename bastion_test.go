package bastion_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/bastion"
)

func TestDefaultBastion(t *testing.T) {
	t.Parallel()

	app := bastion.New()
	e := bastion.Tester(t, app)
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

func TestBastionHelloWorld(t *testing.T) {
	t.Parallel()

	app := bastion.New()
	app.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		render.NewJSON().Send(w, res)
	})

	expected := map[string]interface{}{"message": "world"}

	e := bastion.Tester(t, app)
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestBastionHelloWorldFromFile(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		path string
	}{
		{"from json", "./testdata/options.json"},
		{"from yaml", "./testdata/options.yaml"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app, _ := bastion.FromFile(tc.path)

			assert.Equal(t, ":3000", app.Addr)

			app.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				res := struct {
					Message string `json:"message"`
				}{"world"}
				render.NewJSON().Send(w, res)
			})

			expected := map[string]interface{}{"message": "world"}
			e := bastion.Tester(t, app)
			e.GET("/hello").
				Expect().
				Status(http.StatusOK).
				JSON().Object().Equal(expected)
		})
	}
}

func TestNewRouter(t *testing.T) {
	t.Parallel()

	r := bastion.NewRouter()
	assert.NotNil(t, r)
}

func TestLoadMissingFile(t *testing.T) {
	t.Parallel()

	app, err := bastion.FromFile("./foo.json")
	assert.Nil(t, app)
	assert.Error(t, err, "missing configuration file at ./foo.json")
}

func TestFailUnmarshalFile(t *testing.T) {
	t.Parallel()

	app, err := bastion.FromFile("./testdata/bad_options.json")
	assert.Nil(t, app)
	assert.Error(t, err, "cannot unmarshal configuration file")
}
