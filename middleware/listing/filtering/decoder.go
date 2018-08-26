package filtering

import "net/url"

// Decoder validates filters from url.Values and returns Filters
// with the value found. If not filters were found, it returns nil.
type Decoder struct {
	params   url.Values
	decoders []FilterDecoder
}

// NewDecoder returns a new Decoder instance.
func NewDecoder(params url.Values, decoders ...FilterDecoder) *Decoder {
	return &Decoder{params: params, decoders: decoders}
}

// Decode reads the filter-encoded values from params and stores it
// in the value pointed to by v. If a value is missing from the params
// it'll be filled by their equivalent default value.
func (dec *Decoder) Decode(v *Filtering) error {
	for _, decoder := range dec.decoders {
		filter := decoder.Present(dec.params)
		if filter != nil {
			v.Filters = append(v.Filters, *filter)
		}

		v.Available = append(v.Available, *decoder.WithValues())
	}

	return nil
}
