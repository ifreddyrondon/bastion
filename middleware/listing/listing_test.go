package listing_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalListing(t *testing.T) {
	t.Parallel()
	// sorting
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
	// filtering
	vNew := filtering.NewValue("new", "New")
	vNew.Result = 10
	vUsed := filtering.NewValue("used", "Used")
	vTrue := filtering.NewValue("true", "shared")
	vFalse := filtering.NewValue("false", "private")

	tt := []struct {
		name     string
		l        listing.Listing
		expected string
	}{
		{
			"given a listing with defaults should marshal only paging",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				},
			},
			`{"paging":{"max_allowed_limit":100,"limit":10,"offset":0}}`,
		},
		{
			"given a listing with Paging that includes total should marshal paging with total",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					Total:           1000,
				},
			},
			`{"paging":{"max_allowed_limit":100,"limit":10,"offset":0,"total":1000}}`,
		},
		{
			"given a listing with Paging and Sorting should marshal both",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           20,
					Offset:          10,
					MaxAllowedLimit: 50,
				},
				Sorting: &sorting.Sorting{
					Sort:      &createdDESC,
					Available: []sorting.Sort{createdDESC, createdASC},
				},
			},
			`{"paging":{"max_allowed_limit":50,"limit":20,"offset":10},"sorting":{"sort":{"id":"created_at_desc","description":"Created date descending"},"available":[{"id":"created_at_desc","description":"Created date descending"},{"id":"created_at_asc","description":"Created date ascendant"}]}}`,
		},
		{
			"given a listing with Paging, Sorting and Filtering should marshal all",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           20,
					Offset:          10,
					MaxAllowedLimit: 50,
				},
				Sorting: &sorting.Sorting{
					Sort:      &createdDESC,
					Available: []sorting.Sort{createdDESC, createdASC},
				},
				Filtering: &filtering.Filtering{
					Filters: []filtering.Filter{
						{
							ID:          "condition",
							Description: "test",
							Type:        "text",
							Values:      []filtering.Value{vNew},
						},
					},
					Available: []filtering.Filter{
						{
							ID:          "condition",
							Description: "test",
							Type:        "text",
							Values:      []filtering.Value{vNew, vUsed},
						},
						{
							ID:          "shared",
							Description: "test",
							Type:        "boolean",
							Values:      []filtering.Value{vTrue, vFalse},
						},
					},
				},
			},
			`{"paging":{"max_allowed_limit":50,"limit":20,"offset":10},"sorting":{"sort":{"id":"created_at_desc","description":"Created date descending"},"available":[{"id":"created_at_desc","description":"Created date descending"},{"id":"created_at_asc","description":"Created date ascendant"}]},"filtering":{"filters":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New","result":10}]}],"available":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New","result":10},{"id":"used","description":"Used"}]},{"id":"shared","description":"test","type":"boolean","values":[{"id":"true","description":"shared"},{"id":"false","description":"private"}]}]}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.l.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
