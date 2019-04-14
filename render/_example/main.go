package main

import (
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func main() {
	app := bastion.New()

	app.Get("/data", func(w http.ResponseWriter, req *http.Request) {
		render.Data.Response(w, http.StatusOK, []byte("Some binary data here."))
	})

	app.Get("/text", func(w http.ResponseWriter, req *http.Request) {
		render.Text.Response(w, http.StatusOK, "Plain text here")
	})

	app.Get("/html", func(w http.ResponseWriter, req *http.Request) {
		render.HTML.Response(w, http.StatusOK, "<h1>Hello World</h1>")
	})

	app.Get("/json", func(w http.ResponseWriter, req *http.Request) {
		render.JSON.Response(w, http.StatusOK, map[string]string{"hello": "json"})
	})

	app.Get("/json-ok", func(w http.ResponseWriter, req *http.Request) {
		// with implicit status 200
		render.JSON.Send(w, map[string]string{"hello": "json"})
	})

	app.Get("/json-pretty", func(w http.ResponseWriter, req *http.Request) {
		// with pretty print
		render.
			NewJSON(render.PrettyPrintJSON()).
			Send(w, map[string]string{"hello": "json"})
	})

	app.Get("/xml", func(w http.ResponseWriter, req *http.Request) {
		render.XML.Response(w, http.StatusOK, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Get("/xml-ok", func(w http.ResponseWriter, req *http.Request) {
		// with implicit status 200
		render.XML.Send(w, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Get("/xml-pretty", func(w http.ResponseWriter, req *http.Request) {
		// with pretty print
		render.
			NewXML(render.PrettyPrintXML()).
			Send(w, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Get("/201", func(w http.ResponseWriter, r *http.Request) {
		a := address{"test address", 1, 1}
		render.JSON.Created(w, a)
	})

	app.Get("/400", func(w http.ResponseWriter, r *http.Request) {
		e := errors.New("test")
		render.JSON.BadRequest(w, e)
	})

	app.Get("/500", func(w http.ResponseWriter, r *http.Request) {
		e := errors.New("test")
		render.JSON.InternalServerError(w, e)
	})

	app.Serve()
}
