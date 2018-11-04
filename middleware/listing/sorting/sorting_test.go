package sorting_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalSorting(t *testing.T) {
	t.Parallel()

	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	tt := []struct {
		name     string
		sorting  sorting.Sorting
		expected string
	}{
		{
			"given a empty sorting should marshal without Available",
			sorting.Sorting{},
			`{}`,
		},
		{
			"given a sorting with defaults should marshal defaults",
			sorting.Sorting{
				Sort:      &createdDESC,
				Available: []sorting.Sort{createdDESC},
			},
			`{"sort":{"id":"created_at_desc","description":"Created date descending"},"available":[{"id":"created_at_desc","description":"Created date descending"}]}`,
		},
		{
			"given a sorting with several available should add all to marshal",
			sorting.Sorting{
				Sort:      &createdDESC,
				Available: []sorting.Sort{createdDESC, createdASC},
			},
			`{"sort":{"id":"created_at_desc","description":"Created date descending"},"available":[{"id":"created_at_desc","description":"Created date descending"},{"id":"created_at_asc","description":"Created date ascendant"}]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.sorting.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
