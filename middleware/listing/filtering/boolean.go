package filtering

import "net/url"

const (
	trueID            = "true"
	falseID           = "false"
	booleanFilterType = "boolean"
)

// Boolean validates boolean values and returns a Filter.
// The Filter returned will be true or false.
type Boolean struct {
	id, name   string
	trueValue  Value
	falseValue Value
}

// NewBoolean returns a new Boolean instance.
func NewBoolean(id, name string, trueValName, falseValName string) *Boolean {
	return &Boolean{
		id:         id,
		name:       name,
		trueValue:  NewValue(trueID, trueValName),
		falseValue: NewValue(falseID, falseValName),
	}
}

// Present gets a url and check if a boolean filter is present,
// if it's present validates if its value are true or false.
// Returns a Filter with the applied value or nil is not present.
func (b *Boolean) Present(keys url.Values) *Filter {
	for key, values := range keys {
		if key == b.id {
			v := values[0]
			if v == b.trueValue.ID {
				return NewFilter(b.id, b.name, booleanFilterType, b.trueValue)
			}
			if v == b.falseValue.ID {
				return NewFilter(b.id, b.name, booleanFilterType, b.falseValue)
			}
		}
	}
	return nil
}

// WithValues returns the filter with true and false values.
func (b *Boolean) WithValues() *Filter {
	return NewFilter(b.id, b.name, booleanFilterType, b.trueValue, b.falseValue)
}
