package filtering

import (
	"net/url"

	"github.com/mailru/easyjson/jwriter"
)

// Value is the struct where the posibles filter values should be stored.
type Value struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result int64  `json:"result,omitempty"`
}

// NewValue returns a new Value instance.
func NewValue(id, name string) Value {
	return Value{ID: id, Name: name}
}

// Filter struct that represent a filter
type Filter struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Values []Value `json:"values"`
}

// NewFilter returns a new Filter instance.
func NewFilter(id, name, typef string, values ...Value) *Filter {
	return &Filter{
		ID:     id,
		Name:   name,
		Type:   typef,
		Values: values,
	}
}

// FilterDecoder interface to validate and returns Filter's.
type FilterDecoder interface {
	// Present gets the url params and check if a filter is present within them,
	// if it's present validates if its value is valid.
	// Returns a Filter with the applied value or nil is not present.
	Present(url.Values) *Filter
	// WithValues returns a filter with all their posible values.
	WithValues() *Filter
}

// Filtering allows to filter a collection with the selected Filters
// and their selected values. The Available are all the possible Filters
// with all their possible values.
type Filtering struct {
	Filters   []Filter `json:"filters,omitempty"`
	Available []Filter `json:"available,omitempty"`
}

// MarshalJSON supports json.Marshaler interface
func (v Filtering) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
