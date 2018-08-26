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
	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")
	createdAscSort := sorting.NewSort("created_at_asc", "Created date ascendant")
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
					Sort:      &createdDescSort,
					Available: []sorting.Sort{createdDescSort, createdAscSort},
				},
			},
			`{"paging":{"max_allowed_limit":50,"limit":20,"offset":10},"sorting":{"sort":{"id":"created_at_desc","name":"Created date descending"},"available":[{"id":"created_at_desc","name":"Created date descending"},{"id":"created_at_asc","name":"Created date ascendant"}]}}`,
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
					Sort:      &createdDescSort,
					Available: []sorting.Sort{createdDescSort, createdAscSort},
				},
				Filtering: &filtering.Filtering{
					Filters: []filtering.Filter{
						{
							ID:     "condition",
							Name:   "test",
							Type:   "text",
							Values: []filtering.Value{vNew},
						},
					},
					Available: []filtering.Filter{
						{
							ID:     "condition",
							Name:   "test",
							Type:   "text",
							Values: []filtering.Value{vNew, vUsed},
						},
						{
							ID:     "shared",
							Name:   "test",
							Type:   "boolean",
							Values: []filtering.Value{vTrue, vFalse},
						},
					},
				},
			},
			`{"paging":{"max_allowed_limit":50,"limit":20,"offset":10},"sorting":{"sort":{"id":"created_at_desc","name":"Created date descending"},"available":[{"id":"created_at_desc","name":"Created date descending"},{"id":"created_at_asc","name":"Created date ascendant"}]},"filtering":{"filters":[{"id":"condition","name":"test","type":"text","values":[{"id":"new","name":"New","result":10}]}],"available":[{"id":"condition","name":"test","type":"text","values":[{"id":"new","name":"New","result":10},{"id":"used","name":"Used"}]},{"id":"shared","name":"test","type":"boolean","values":[{"id":"true","name":"shared"},{"id":"false","name":"private"}]}]}}`,
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
