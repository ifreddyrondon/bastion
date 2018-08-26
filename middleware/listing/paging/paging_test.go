package paging_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalPaging(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		paging   paging.Paging
		expected string
	}{
		{
			"given a paging with total marshal with total",
			paging.Paging{
				Total:           300,
				Offset:          100,
				Limit:           20,
				MaxAllowedLimit: 100,
			},
			`{"max_allowed_limit":100,"limit":20,"offset":100,"total":300}`,
		},
		{
			"given a paging with total equals to 0 marshal should omit total",
			paging.Paging{
				Total:           0,
				Offset:          100,
				Limit:           20,
				MaxAllowedLimit: 100,
			},
			`{"max_allowed_limit":100,"limit":20,"offset":100}`,
		},
		{
			"given a paging without total marshal should omit total",
			paging.Paging{
				Offset:          100,
				Limit:           20,
				MaxAllowedLimit: 100,
			},
			`{"max_allowed_limit":100,"limit":20,"offset":100}`,
		},
		{
			"given a paging with Offset, Limit or MaxAllowedLimit equals to 0 marshal should always include them",
			paging.Paging{},
			`{"max_allowed_limit":0,"limit":0,"offset":0}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.paging.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
