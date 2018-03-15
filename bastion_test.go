package bastion_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func TestDefaultBastion(t *testing.T) {
	app := bastion.New(bastion.Options{})
	e := bastion.Tester(t, app)
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		Text().Equal("pong")
}

func TestBastionHelloWorld(t *testing.T) {
	app := bastion.New(bastion.Options{})
	app.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		res := struct {
			Message string `json:"message"`
		}{"world"}
		err := json.NewRender(w).Send(res)
		assert.Nil(t, err)
	})

	expected := map[string]interface{}{"message": "world"}

	e := bastion.Tester(t, app)
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(expected)
}

func TestBastionHelloWorldFromFile(t *testing.T) {
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

			assert.Equal(t, "/api/", app.APIBasepath)
			assert.Equal(t, ":3000", app.Addr)
			assert.True(t, app.Debug)

			app.APIRouter.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				res := struct {
					Message string `json:"message"`
				}{"world"}
				err := json.NewRender(w).Send(res)
				assert.Nil(t, err)
			})

			expected := map[string]interface{}{"message": "world"}
			e := bastion.Tester(t, app)
			e.GET("/api/hello").
				Expect().
				Status(http.StatusOK).
				JSON().Object().Equal(expected)
		})
	}
}

func TestNewRouter(t *testing.T) {
	r := bastion.NewRouter()
	assert.NotNil(t, r)
}

func TestBastionFromPartialYAMLFile(t *testing.T) {
	app, _ := bastion.FromFile("./testdata/partial_options.yaml")
	assert.Equal(t, "/api/", app.APIBasepath)
	assert.Equal(t, "127.0.0.1:8080", app.Addr)
	assert.False(t, app.Debug)
}

func TestLoadMissingFile(t *testing.T) {
	app, err := bastion.FromFile("./foo.json")
	assert.Nil(t, app)
	assert.Error(t, err, "missing configuration file at ./foo.json")
}

func TestFailUnmarshalFile(t *testing.T) {
	app, err := bastion.FromFile("./testdata/bad_options.json")
	assert.Nil(t, app)
	assert.Error(t, err, "cannot unmarshal configuration file")
}
