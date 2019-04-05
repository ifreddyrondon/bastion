package filtering_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
)

func TestTextPresentOK(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")

	tt := []struct {
		name     string
		text     *filtering.Text
		params   url.Values
		expected *filtering.Filter
	}{
		{
			"should return new value when param with condition with new value",
			filtering.NewText("condition", "test", vNew),
			map[string][]string{"condition": {"new"}},
			&filtering.Filter{
				ID:          "condition",
				Description: "test",
				Type:        "text",
				Values:      []filtering.Value{filtering.NewValue("new", "New")},
			},
		},
		{
			"should return used value when param with condition with used value",
			filtering.NewText("condition", "test", vUsed),
			map[string][]string{"condition": {"used"}},
			&filtering.Filter{
				ID:          "condition",
				Description: "test",
				Type:        "text",
				Values:      []filtering.Value{filtering.NewValue("used", "Used")},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.text.Present(tc.params)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Description, result.Description)
			assert.Equal(t, tc.expected.Type, result.Type)
			assert.Equal(t, len(tc.expected.Values), len(result.Values))
			assert.Equal(t, tc.expected.Values[0].ID, result.Values[0].ID)
			assert.Equal(t, tc.expected.Values[0].Description, result.Values[0].Description)
		})
	}
}

func TestTextPresentFails(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	tt := []struct {
		name   string
		text   *filtering.Text
		params url.Values
	}{
		{
			"should return nil when not value found",
			filtering.NewText("condition", "test", vNew),
			map[string][]string{"condition": {"abc"}},
		},
		{
			"should return nil when not params found",
			filtering.NewText("condition", "test", vNew),
			map[string][]string{"foo": {"abc"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.text.Present(tc.params)
			assert.Nil(t, result)
		})
	}
}

func TestTextWithValues(t *testing.T) {
	t.Parallel()

	values := []filtering.Value{
		filtering.NewValue("new", "New"),
		filtering.NewValue("used", "Used"),
	}
	text := filtering.NewText("condition", "test", values...)
	expected := &filtering.Filter{
		ID:          "condition",
		Description: "test",
		Type:        "text",
		Values:      values,
	}
	result := text.WithValues()
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Description, result.Description)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, len(expected.Values), len(result.Values))
	for i, v := range result.Values {
		assert.Equal(t, expected.Values[i].ID, v.ID)
		assert.Equal(t, expected.Values[i].Description, v.Description)
	}
}
