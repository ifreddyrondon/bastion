package binder_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/binder"
)

func TestNewXMLFromReq(t *testing.T) {
	t.Parallel()
	body := `
<?xml version="1.0" encoding="UTF-8"?>
<address>
   <address>la comarca</address>
   <lat>123</lat>
</address>
`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewXML().FromReq(req, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestXMLFromReqMissingBody(t *testing.T) {
	t.Parallel()
	var a address
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := binder.XML.FromReq(req, &a)
	assert.EqualError(t, err, "invalid request, body not present")
}

func TestXMLFromSrc(t *testing.T) {
	t.Parallel()
	src := []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<address>
   <address>la comarca</address>
   <lat>123</lat>
</address>
`)
	var a address
	err := binder.XML.FromSrc(src, &a)
	assert.Nil(t, err)
	assert.Equal(t, "la comarca", a.Address)
	assert.Equal(t, 123.0, a.Lat)
	assert.Equal(t, 0.0, a.Lng)
}

func TestXMLDecodeWithXMLDecodingErrMsg(t *testing.T) {
	t.Parallel()
	body := `{`
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	var a address
	err := binder.NewXML(binder.XMLDecodingErrMsg("test")).FromReq(req, &a)
	assert.EqualError(t, err, "test")
}

func TestXMLValidate(t *testing.T) {
	t.Parallel()
	body := []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<address>
   <address>la comarca</address>
   <lat>-1</lat>
</address>
`)
	var a address
	err := binder.XML.FromSrc(body, &a)
	assert.EqualError(t, err, "address lat can't be lower than 0")
}
