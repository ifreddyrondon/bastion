package main

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"
)

type address struct {
	Address *string `json:"address" xml:"address" yaml:"address"`
	Lat     float64 `json:"lat" xml:"lat" yaml:"lat"`
	Lng     float64 `json:"lng" xml:"lng" yaml:"lng"`
}

func (a *address) Validate() error {
	if a.Address == nil || *a.Address == "" {
		return errors.New("missing address field")
	}
	return nil
}

func main() {
	app := bastion.New()
	app.Post("/decode-json", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.JSON.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Post("/decode-xml", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.XML.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Post("/decode-yaml", func(w http.ResponseWriter, r *http.Request) {
		var a address
		if err := binder.YAML.FromReq(r, &a); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}
		render.JSON.Send(w, a)
	})
	app.Serve()
}
