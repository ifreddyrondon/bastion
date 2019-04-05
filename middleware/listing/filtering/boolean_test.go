package filtering_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
)

func TestBooleanPresentOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		boolean  *filtering.Boolean
		params   url.Values
		expected *filtering.Filter
	}{
		{
			"should return true value when url params contains a true value",
			filtering.NewBoolean("shared", "test", "shared", "private"),
			map[string][]string{"shared": {"true"}},
			&filtering.Filter{
				ID:          "shared",
				Description: "test",
				Type:        "boolean",
				Values:      []filtering.Value{filtering.NewValue("true", "shared")},
			},
		},
		{
			"should return false value when url param contains a false value",
			filtering.NewBoolean("shared", "test", "shared", "private"),
			map[string][]string{"shared": {"false"}},
			&filtering.Filter{
				ID:          "shared",
				Description: "test",
				Type:        "boolean",
				Values:      []filtering.Value{filtering.NewValue("false", "private")},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.boolean.Present(tc.params)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.Type, result.Type)
			assert.Equal(t, len(tc.expected.Values), len(result.Values))
			assert.Equal(t, tc.expected.Values[0].ID, result.Values[0].ID)
			assert.Equal(t, tc.expected.Values[0].Description, result.Values[0].Description)
		})
	}
}

func TestBooleanPresentFails(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		boolean *filtering.Boolean
		params  url.Values
	}{
		{
			"should return nil when not value found",
			filtering.NewBoolean("shared", "test", "shared", "private"),
			map[string][]string{"shared": {"abc"}},
		},
		{
			"should return nil when not params found",
			filtering.NewBoolean("shared", "test", "shared", "private"),
			map[string][]string{"foo": {"abc"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.boolean.Present(tc.params)
			assert.Nil(t, result)
		})
	}
}

func TestBooleanWithValues(t *testing.T) {
	t.Parallel()

	boolean := filtering.NewBoolean("shared", "test", "shared", "private")
	expected := &filtering.Filter{
		ID:          "shared",
		Description: "test",
		Type:        "boolean",
		Values: []filtering.Value{
			filtering.NewValue("true", "shared"),
			filtering.NewValue("false", "private"),
		},
	}
	result := boolean.WithValues()
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Description, result.Description)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, len(expected.Values), len(result.Values))
	assert.Equal(t, expected.Values[0].ID, result.Values[0].ID)
	assert.Equal(t, expected.Values[0].Description, result.Values[0].Description)
	assert.Equal(t, expected.Values[1].ID, result.Values[1].ID)
	assert.Equal(t, expected.Values[1].Description, result.Values[1].Description)
}
