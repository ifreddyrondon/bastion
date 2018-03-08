package json_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ifreddyrondon/bastion/reader/json"
	"github.com/stretchr/testify/assert"
)

type readerCloser struct {
	io.Reader
}

func (t readerCloser) Close() error { return nil }

type Location struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func getBody(address string, lat, lng float64) io.ReadCloser {
	payload := strings.NewReader(fmt.Sprintf(`{"address":"%v","lat":%v,"lng":%v}`, address, lat, lng))
	return &readerCloser{payload}
}

func TestReadJSONWithDefinedStrut(t *testing.T) {
	var addr Location
	expected := Location{"jorge matte gormaz", 1, 1}

	body := getBody(expected.Address, expected.Lat, expected.Lng)
	err := json.NewReader(body).Read(&addr)

	assert.Nil(t, err)
	assert.Equal(t, expected.Address, addr.Address)
	assert.Equal(t, expected.Lat, addr.Lat)
	assert.Equal(t, expected.Lng, addr.Lng)
}

func TestReadJSONWithMap(t *testing.T) {
	result := make(map[string]interface{})
	expected := Location{"jorge matte gormaz", 1, 1}

	body := getBody(expected.Address, expected.Lat, expected.Lng)
	err := json.NewReader(body).Read(&result)

	assert.Nil(t, err)
	assert.Equal(t, expected.Address, result["address"])
	assert.Equal(t, expected.Lat, result["lat"])
	assert.Equal(t, expected.Lng, result["lng"])
}

func TestReadJSONError(t *testing.T) {
	container := make(map[string]interface{})
	body := &readerCloser{strings.NewReader("`")}
	expectedErr := "invalid character '`' looking for beginning of value"
	err := json.NewReader(body).Read(&container)

	assert.EqualError(t, err, expectedErr)
}
