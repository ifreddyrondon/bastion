package paging

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	// DefaultMaxAllowedLimit determines the maximum value that can take limit when decode.
	DefaultMaxAllowedLimit = 100
	// DefaultLimit default value that take limit when not limit param present in query params.
	DefaultLimit = 10
	// DefaultOffset default value that take offset when not offset param present in query params.
	DefaultOffset = 0
)

var (
	// ErrInvalidOffsetValueNotANumber expected error when fails parsing the offset value to int.
	ErrInvalidOffsetValueNotANumber = errors.New("invalid offset value, must be a number")
	// ErrInvalidOffsetValueLessThanZero expected error when offset value is less than zero.
	ErrInvalidOffsetValueLessThanZero = errors.New("invalid offset value, must be greater than zero")
	// ErrInvalidLimitValueNotANumber expected error when fails parse limit value to int
	ErrInvalidLimitValueNotANumber = errors.New("invalid limit value, must be a number")
	// ErrInvalidLimitValueLessThanZero expected error when limit value is less than zero.
	ErrInvalidLimitValueLessThanZero = errors.New("invalid limit value, must be greater than zero")
)

// Option allows to modify the defaults decode values.
type Option func(*Decoder)

// A Decoder reads and decodes Paging values from url.Values.
type Decoder struct {
	params url.Values

	maxAllowedLimit int
	limit           int
	offset          int64
}

// Limit set the paging limit default.
func Limit(limit int) Option {
	return func(dec *Decoder) {
		dec.limit = limit
	}
}

// MaxAllowedLimit set the max allowed limit default.
func MaxAllowedLimit(maxAllowed int) Option {
	return func(dec *Decoder) {
		dec.maxAllowedLimit = maxAllowed
	}
}

// NewDecoder returns a new decoder that reads from params.
func NewDecoder(params url.Values, opts ...Option) *Decoder {
	d := &Decoder{
		params:          params,
		maxAllowedLimit: DefaultMaxAllowedLimit,
		limit:           DefaultLimit,
		offset:          DefaultOffset,
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

// Decode reads the next Paging-encoded value from
// params and stores it in the value pointed to by v.
func (dec *Decoder) Decode(v *Paging) error {
	dec.fillDefaults(v)
	offsetStr, ok := dec.params["offset"]
	if ok {
		off, err := strconv.ParseInt(offsetStr[0], 10, 64)
		if err != nil {
			return ErrInvalidOffsetValueNotANumber
		}
		if off < 0 {
			return ErrInvalidOffsetValueLessThanZero
		}
		v.Offset = off
	}
	limitStr, ok := dec.params["limit"]
	if ok {
		l, err := strconv.Atoi(limitStr[0])
		if err != nil {
			return ErrInvalidLimitValueNotANumber
		}
		if l < 0 {
			return ErrInvalidLimitValueLessThanZero
		}
		if l > dec.maxAllowedLimit {
			l = dec.maxAllowedLimit
		}
		v.Limit = l
	}

	return nil
}

func (dec *Decoder) fillDefaults(p *Paging) {
	p.Offset = dec.offset
	p.Limit = dec.limit
	p.MaxAllowedLimit = dec.maxAllowedLimit
}
