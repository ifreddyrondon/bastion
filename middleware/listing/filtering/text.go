package filtering

import "net/url"

const textFilterType = "text"

// Text validates text values and returns a Filter with the
// value found. If not filters were found, it returns nil.
type Text struct {
	id, description string
	values          []Value
}

// NewText returns a new Text instance.
func NewText(id, description string, values ...Value) *Text {
	return &Text{
		id:          id,
		description: description,
		values:      values,
	}
}

// Present gets the url params and check if a text filter is present,
// if it's present validates its value meets one of filter values options.
// Returns a Filter with the applied value or nil is not present.
func (b *Text) Present(keys url.Values) *Filter {
	for key, values := range keys {
		if key == b.id {
			v := checkValues(b.values, values[0])
			if v != nil {
				return NewFilter(b.id, b.description, textFilterType, *v)
			}
		}
	}
	return nil
}

func checkValues(available []Value, paramVal string) *Value {
	for _, v := range available {
		if v.ID == paramVal {
			return &v
		}
	}
	return nil
}

// WithValues returns the filter with all their values.
func (b *Text) WithValues() *Filter {
	return NewFilter(b.id, b.description, textFilterType, b.values...)
}
