package render_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/render"
)

type addressXML struct {
	Address string  `xml:"address,omitempty"`
	Lat     float64 `xml:"lat,omitempty"`
	Lng     float64 `xml:"lng,omitempty"`
}

func TestXMLResponse(t *testing.T) {
	t.Parallel()

	a := &addressXML{"test address", 1, 1}

	tt := []struct {
		name     string
		opts     []func(*render.XMLRenderer)
		a        *addressXML
		expected string
	}{
		{
			"marshal without indent",
			[]func(*render.XMLRenderer){},
			a,
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<addressXML><address>test address</address><lat>1</lat><lng>1</lng></addressXML>",
		},
		{
			"marshal with indent (pretty print)",
			[]func(*render.XMLRenderer){render.PrettyPrintXML()},
			a,
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n  <addressXML>\n      <address>test address</address>\n      <lat>1</lat>\n      <lng>1</lng>\n  </addressXML>",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			render.NewXML(tc.opts...).Response(rr, http.StatusOK, tc.a)
			httpexpect.NewResponse(t, rr.Result()).
				Status(http.StatusOK).
				Body().
				Equal(tc.expected)
		})
	}
}

func TestXMLResponseError(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	render.XML.Response(rr, http.StatusOK, make(chan int))
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusInternalServerError).
		Text().
		Equal("xml: unsupported type: chan int\n")
}

func TestXMLSend(t *testing.T) {
	t.Parallel()

	a := addressXML{"test address", 1, 1}
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<addressXML><address>test address</address><lat>1</lat><lng>1</lng></addressXML>"

	rr := httptest.NewRecorder()
	render.XML.Send(rr, &a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusOK).
		Body().
		Equal(expected)
}

func TestXMLCreated(t *testing.T) {
	t.Parallel()

	a := addressXML{"test address", 1, 1}
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<addressXML><address>test address</address><lat>1</lat><lng>1</lng></addressXML>"

	rr := httptest.NewRecorder()
	render.XML.Created(rr, &a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusCreated).
		Body().
		Equal(expected)
}

func TestXMLNoContent(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	render.XML.NoContent(rr)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNoContent).NoContent()
}

func TestXMLBadRequest(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<HTTPError message=\"test\" error=\"Bad Request\" status=\"400\"></HTTPError>"

	rr := httptest.NewRecorder()
	render.XML.BadRequest(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusBadRequest).
		Body().
		Equal(expected)
}

func TestXMLNotFound(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<HTTPError message=\"test\" error=\"Not Found\" status=\"404\"></HTTPError>"

	rr := httptest.NewRecorder()
	render.XML.NotFound(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNotFound).
		Body().
		Equal(expected)
}

func TestXMLMethodNotAllowed(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<HTTPError message=\"test\" error=\"Method Not Allowed\" status=\"405\"></HTTPError>"

	rr := httptest.NewRecorder()
	render.XML.MethodNotAllowed(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusMethodNotAllowed).
		Body().
		Equal(expected)
}

func TestXMLInternalServerError(t *testing.T) {
	t.Parallel()

	e := errors.New("test")
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<HTTPError message=\"test\" error=\"Internal Server Error\" status=\"500\"></HTTPError>"

	rr := httptest.NewRecorder()
	render.XML.InternalServerError(rr, e)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusInternalServerError).
		Body().
		Equal(expected)
}
