package binder_test

import "github.com/pkg/errors"

type address struct {
	Address string  `json:"address" xml:"address" yaml:"address"`
	Lat     float64 `json:"lat" xml:"lat" yaml:"lat"`
	Lng     float64 `json:"lng" xml:"lng" yaml:"lng"`
}

func (n *address) Validate() error {
	if n.Lat < 0 {
		return errors.New("address lat can't be lower than 0")
	}
	return nil
}
