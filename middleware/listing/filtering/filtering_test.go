package filtering_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
)

func TestMarshalSorting(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	vNewWithResult := filtering.NewValue("new", "New")
	vNewWithResult.Result = 10
	vUsed := filtering.NewValue("used", "Used")
	vTrue := filtering.NewValue("true", "shared")
	vFalse := filtering.NewValue("false", "private")

	tt := []struct {
		name     string
		f        filtering.Filtering
		expected string
	}{
		{
			"given a empty Filtering should marshal empty object {}",
			filtering.Filtering{},
			`{}`,
		},
		{
			"given a filtering with non aplied Filters and one available filter should marshal that filter as available",
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					{
						ID:          "condition",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{vNew, vUsed},
					},
				},
			},
			`{"available":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New"},{"id":"used","description":"Used"}]}]}`,
		},
		{
			"given a filtering with one aplied Filter with result should marshal the applied filter and the availables",
			filtering.Filtering{
				Filters: []filtering.Filter{
					{
						ID:          "condition",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{vNewWithResult},
					},
				},
				Available: []filtering.Filter{
					{
						ID:          "condition",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{vNewWithResult, vUsed},
					},
				},
			},
			`{"filters":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New","result":10}]}],"available":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New","result":10},{"id":"used","description":"Used"}]}]}`,
		},
		{
			"given a filtering with one aplied Filters and some available filters should marshal them as available and the applied filter",
			filtering.Filtering{
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
			`{"filters":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New"}]}],"available":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New"},{"id":"used","description":"Used"}]},{"id":"shared","description":"test","type":"boolean","values":[{"id":"true","description":"shared"},{"id":"false","description":"private"}]}]}`,
		},
		{
			"given a filtering with all aplied Filters and available filters should marshal them as available and the all filters",
			filtering.Filtering{
				Filters: []filtering.Filter{
					{
						ID:          "condition",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{vNew},
					},
					{
						ID:          "shared",
						Description: "test",
						Type:        "boolean",
						Values:      []filtering.Value{vFalse},
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
			`{"filters":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New"}]},{"id":"shared","description":"test","type":"boolean","values":[{"id":"false","description":"private"}]}],"available":[{"id":"condition","description":"test","type":"text","values":[{"id":"new","description":"New"},{"id":"used","description":"Used"}]},{"id":"shared","description":"test","type":"boolean","values":[{"id":"true","description":"shared"},{"id":"false","description":"private"}]}]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.f.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
