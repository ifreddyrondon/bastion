# Render

Easily rendering JSON, XML, binary data, and HTML templates responses 

## Usage
It can be used with pretty much any web framework providing you can access the `http.ResponseWriter` from your handler.
The rendering functions simply wraps Go's existing functionality for marshaling and rendering data.

```go
package main

import (
	"encoding/xml"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
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

	app.Get("/xml", func(w http.ResponseWriter, req *http.Request) {
		render.XML.Response(w, http.StatusOK, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Get("/xml-ok", func(w http.ResponseWriter, req *http.Request) {
		// with implicit status 200
		render.XML.Send(w, ExampleXML{One: "hello", Two: "xml"})
	})

	app.Serve()
}
```

### Renderer

```go
// Renderer interface for managing response payloads.
type Renderer interface {
	// Response encoded responses in the ResponseWriter with the HTTP status code.
	Response(w http.ResponseWriter, code int, response interface{})
}
```
Implementations: 

- **render.Data** response []byte with application/octet-stream Content-Type.
- **render.Text** response strings with text/plain Content-Type.
- **render.HTML** response strings with text/html Content-Type.
- **render.JSON** response strings with application/json Content-Type.
- **render.XML** response strings with application/xml Content-Type.

### APIRenderer

APIRenderer are convenient methods for api responses.

```go
// APIRenderer interface for managing API response payloads.
type APIRenderer interface {
	Renderer
	OKRenderer
	ClientErrRenderer
	ServerErrRenderer
}

// OKRenderer interface for managing success API response payloads.
type OKRenderer interface {
	Send(w http.ResponseWriter, response interface{})
	Created(w http.ResponseWriter, response interface{})
	NoContent(w http.ResponseWriter)
}

// ClientErrRenderer interface for managing API responses when client error.
type ClientErrRenderer interface {
	BadRequest(w http.ResponseWriter, err error)
	NotFound(w http.ResponseWriter, err error)
	MethodNotAllowed(w http.ResponseWriter, err error)
}

// ServerErrRenderer interface for managing API responses when server error.
type ServerErrRenderer interface {
	InternalServerError(w http.ResponseWriter, err error)
}
```

Implementations:

- **render.JSON** response strings with text/html Content-Type.
- **render.XML** response strings with text/html Content-Type.

```go
package main

import (
	"net/http"
	"errors"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func h201(w http.ResponseWriter, r *http.Request) {
	a := address{"test address", 1, 1}
	render.JSON.Created(w, a)
}

func h400(w http.ResponseWriter, r *http.Request) {
	e := errors.New("test")
	render.JSON.BadRequest(w, e)
}

func h500(w http.ResponseWriter, r *http.Request) {
	e := errors.New("test")
	render.JSON.InternalServerError(w, e)
}

func main() {
	app := bastion.New()
	app.Get("/201", h201)
	app.Get("/400", h400)
	app.Get("/500", h500)
	app.Serve()
}
```

With Options:

```go
package main

import (
	"encoding/xml"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type ExampleXML struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	render.
			NewJSON(render.PrettyPrintJSON()).
			Send(w, map[string]string{"hello": "json with options"})
}

func xmlHandler(w http.ResponseWriter, r *http.Request) {
	render.
			NewXML(render.PrettyPrintXML()).
			Send(w, ExampleXML{One: "hello", Two: "xml"})
}

func main() {
	app := bastion.New()
	app.Get("/json", jsonHandler)
	app.Get("/xml", xmlHandler)
	app.Serve()
}
```

[**E.g.**](https://github.com/ifreddyrondon/bastion/blob/master/render/_example/main.go)



