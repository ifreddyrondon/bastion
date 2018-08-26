package sorting_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")
	createdAscSort := sorting.NewSort("created_at_asc", "Created date ascendant")

	tt := []struct {
		name      string
		urlParams url.Values
		criteria  []sorting.Sort
		result    sorting.Sorting
	}{
		{
			"given none query params and non criteria should decode empty Sorting",
			map[string][]string{},
			[]sorting.Sort{},
			sorting.Sorting{},
		},
		{
			"given non sort query params present and one sort criteria",
			map[string][]string{},
			[]sorting.Sort{createdDescSort},
			sorting.Sorting{
				Sort:      &createdDescSort,
				Available: []sorting.Sort{createdDescSort},
			},
		},
		{
			"given non sort query params present and one some sort criteria",
			map[string][]string{},
			[]sorting.Sort{createdDescSort, createdAscSort},
			sorting.Sorting{
				Sort:      &createdDescSort,
				Available: []sorting.Sort{createdDescSort, createdAscSort},
			},
		},
		{
			"given created_at_desc sort query params present and one some sort criteria",
			map[string][]string{"sort": {"created_at_desc"}},
			[]sorting.Sort{createdDescSort, createdAscSort},
			sorting.Sorting{
				Sort:      &createdDescSort,
				Available: []sorting.Sort{createdDescSort, createdAscSort},
			},
		},
		{
			"given created_at_desc sort query params present and one some sort criteria",
			map[string][]string{"sort": {"created_at_asc"}},
			[]sorting.Sort{createdDescSort, createdAscSort},
			sorting.Sorting{
				Sort:      &createdAscSort,
				Available: []sorting.Sort{createdDescSort, createdAscSort},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s sorting.Sorting
			err := sorting.NewDecoder(tc.urlParams, tc.criteria...).Decode(&s)
			assert.Nil(t, err)
			assert.Equal(t, tc.result.Sort, s.Sort)
			assert.Equal(t, len(tc.result.Available), len(s.Available))
		})
	}
}

func TestSortingDecodeBad(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams url.Values
		criteria  []sorting.Sort
		err       string
	}{
		{
			"given a sort query when non sorting criteria",
			map[string][]string{"sort": {"foo_desc"}},
			[]sorting.Sort{},
			"there's no order criteria with the id foo_desc",
		},
		{
			"given a sort query when none match sorting criteria",
			map[string][]string{"sort": {"foo_desc"}},
			[]sorting.Sort{
				sorting.NewSort("created_at_desc", "Created date descending"),
			},
			"there's no order criteria with the id foo_desc",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s sorting.Sorting
			err := sorting.NewDecoder(tc.urlParams, tc.criteria...).Decode(&s)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
