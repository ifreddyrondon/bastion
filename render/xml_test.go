package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/bastion/render"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
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
		opts     []func(*render.XML)
		a        *addressXML
		expected string
	}{
		{
			"marshal without indent",
			[]func(*render.XML){},
			a,
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<addressXML><address>test address</address><lat>1</lat><lng>1</lng></addressXML>",
		},
		{
			"marshal with indent (pretty print)",
			[]func(*render.XML){render.PrettyPrintXML()},
			a,
			"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n  <addressXML>\n      <address>test address</address>\n      <lat>1</lat>\n      <lng>1</lng>\n  </addressXML>",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			render.NewXML(rr, tc.opts...).Response(http.StatusOK, tc.a)
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
	render.NewXML(rr).Response(http.StatusOK, make(chan int))
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
	render.NewXML(rr).Send(&a)
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
	render.NewXML(rr).Created(&a)
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusCreated).
		Body().
		Equal(expected)
}

func TestXMLNoContent(t *testing.T) {
	t.Parallel()

	rr := httptest.NewRecorder()
	render.NewXML(rr).NoContent()
	httpexpect.NewResponse(t, rr.Result()).
		Status(http.StatusNoContent).NoContent()
}
