package gobastion_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ifreddyrondon/gobastion"
)

type testerReaderCloser struct {
	io.Reader
}

func (t testerReaderCloser) Close() error { return nil }

func addressToBytes(address string, lat, lng float64) []byte {
	res := fmt.Sprintf(`{"address":"%v", "lat":%v, "lng":%v}`, address, lat, lng)
	return []byte(res)
}

func TestReadJSONWithDefinedStrut(t *testing.T) {
	reader := gobastion.JsonReader{}
	container := struct {
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}{}
	expected := struct {
		address  string
		lat, lng float64
	}{"jorge matte gormaz", 1, 1}

	payload := addressToBytes(expected.address, expected.lat, expected.lng)
	input := testerReaderCloser{bytes.NewBuffer(payload)}
	err := reader.Read(&input, &container)

	if err != nil {
		t.Fatalf("Expected err to be nil. Got '%v'", err)
	}
	if expected.address != container.Address {
		t.Errorf("Expected address to be '%v'. Got '%v'", expected.address, container.Address)
	}
	if expected.lat != container.Lat {
		t.Errorf("Expected lat to be '%v'. Got '%v'", expected.lat, container.Lat)
	}
	if expected.lng != container.Lng {
		t.Errorf("Expected lng to be '%v'. Got '%v'", expected.lng, container.Lng)
	}
}

func TestReadJSONWithMap(t *testing.T) {
	reader := gobastion.JsonReader{}
	container := make(map[string]interface{})
	expected := struct {
		address  string
		lat, lng float64
	}{"jorge matte gormaz", 1, 1}

	payload := addressToBytes(expected.address, expected.lat, expected.lng)
	input := testerReaderCloser{bytes.NewBuffer(payload)}
	err := reader.Read(&input, &container)

	if err != nil {
		t.Fatalf("Expected err to be nil. Got '%v'", err)
	}
	if expected.address != container["address"] {
		t.Errorf("Expected address to be '%v'. Got '%v'", expected.address, container["address"])
	}
	if expected.lat != container["lat"] {
		t.Errorf("Expected lat to be '%v'. Got '%v'", expected.lat, container["lat"])
	}
	if expected.lng != container["lng"] {
		t.Errorf("Expected lng to be '%v'. Got '%v'", expected.lng, container["lng"])
	}
}

func TestReadJSONError(t *testing.T) {
	reader := gobastion.JsonReader{}
	container := make(map[string]interface{})
	input := testerReaderCloser{strings.NewReader("`")}
	expectedErr := "invalid character '`' looking for beginning of value"
	err := reader.Read(&input, &container)

	if expectedErr != err.Error() {
		t.Fatalf("Expected err to be '%v'. Got '%v'", expectedErr, err.Error())
	}
}
