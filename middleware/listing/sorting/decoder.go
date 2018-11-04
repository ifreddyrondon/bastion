package sorting

import (
	"fmt"
	"net/url"
)

const errSortKeyNotAvailable = "there's no order criteria with the id %v"

// A Decoder reads and decodes Sorting values from url.Values.
type Decoder struct {
	params url.Values

	defaultCriteria *Sort
	criteria        []Sort
}

// NewDecoder returns a new decoder that reads from query params.
//
// The first sort criteria if present will be the default Sort
// when decode url query params and not params present.
func NewDecoder(params url.Values, criteria ...Sort) *Decoder {
	d := &Decoder{params: params, criteria: criteria}

	if len(criteria) > 0 {
		d.defaultCriteria = &criteria[0]
	}

	return d
}

// Decode reads the sort-encoded value from params and stores it
// in the value pointed to by v. If a value is missing from the params
// it'll be filled by their equivalent default value.
func (dec *Decoder) Decode(v *Sorting) error {
	dec.fillDefaults(v)

	sortStr, ok := dec.params["sort"]
	if ok {
		sort := paramsInAvailable(sortStr[0], dec.criteria)
		if sort == nil {
			return fmt.Errorf(errSortKeyNotAvailable, sortStr[0])
		}
		v.Sort = sort
	}
	return nil
}

func paramsInAvailable(sortKey string, available []Sort) *Sort {
	for _, sort := range available {
		if sortKey == sort.ID {
			return &sort
		}
	}

	return nil
}

func (dec *Decoder) fillDefaults(s *Sorting) {
	s.Available = dec.criteria
	// fix: avoid override defaultCriteria values when change Sorting
	if dec.defaultCriteria != nil {
		sort := *dec.defaultCriteria
		s.Sort = &sort
	}
}
