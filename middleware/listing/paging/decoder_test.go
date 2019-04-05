package paging_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams url.Values
		opts      []paging.Option
		result    paging.Paging
	}{
		{
			"given none query params and non options should decode defaults",
			map[string][]string{},
			[]paging.Option{},
			func() paging.Paging {
				return paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given none query params and limit option should decode defaults with new limit",
			map[string][]string{},
			[]paging.Option{paging.Limit(50)},
			func() paging.Paging {
				return paging.Paging{
					Limit:           50,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given a limit param and non options should decode with limit and defaults",
			map[string][]string{"limit": {"1"}},
			[]paging.Option{},
			func() paging.Paging {
				return paging.Paging{
					Limit:           1,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given a offset param and non options should decode with offset and defaults",
			map[string][]string{"offset": {"1"}},
			[]paging.Option{},
			func() paging.Paging {
				return paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          1,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given offset and limit params and non options should decode with offset, limit and defaults",
			map[string][]string{"offset": {"1"}, "limit": {"1"}},
			[]paging.Option{},
			func() paging.Paging {
				return paging.Paging{
					Limit:           1,
					Offset:          1,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given offset and limit when limit > maxAllowed default and none option should decode with offset and limit set to default",
			map[string][]string{"offset": {"1"}, "limit": {"101"}},
			[]paging.Option{},
			func() paging.Paging {
				return paging.Paging{
					Limit:           100,
					Offset:          1,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				}
			}(),
		},
		{
			"given offset and limit when limit > maxAllowed default and maxAllowed option should decode with offset and limit upper the default",
			map[string][]string{"offset": {"1"}, "limit": {"105"}},
			[]paging.Option{paging.MaxAllowedLimit(110)},
			func() paging.Paging {
				return paging.Paging{
					Limit:           105,
					Offset:          1,
					MaxAllowedLimit: 110,
				}
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var p paging.Paging
			err := paging.NewDecoder(tc.urlParams, tc.opts...).Decode(&p)
			assert.Nil(t, err)
			assert.Equal(t, p.Limit, tc.result.Limit)
			assert.Equal(t, p.Offset, tc.result.Offset)
			assert.Equal(t, p.MaxAllowedLimit, tc.result.MaxAllowedLimit)
		})
	}
}

func TestDecodeFails(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams url.Values
		err       string
	}{
		{
			"given a not number limit param should return an error",
			map[string][]string{"limit": {"a"}},
			"invalid limit value, must be a number",
		},
		{
			"given a limit < 0 param should return an error",
			map[string][]string{"limit": {"-1"}},
			"invalid limit value, must be greater than zero",
		},
		{
			"given a not number offset param should return an error",
			map[string][]string{"offset": {"a"}},
			"invalid offset value, must be a number",
		},
		{
			"given a offset < 0 param should return an error",
			map[string][]string{"offset": {"-1"}},
			"invalid offset value, must be greater than zero",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var p paging.Paging
			err := paging.NewDecoder(tc.urlParams).Decode(&p)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
