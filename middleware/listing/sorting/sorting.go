package sorting

import "github.com/mailru/easyjson/jwriter"

// Sort criteria.
type Sort struct {
	ID          string `json:"id"`
	Value       string `json:"-"`
	Description string `json:"description"`
}

// MarshalJSON supports json.Marshaler interface
func (v Sorting) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1afec5d2EncodeGithubComIfreddyrondonCaptureAppListingSorting(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// NewSort returns a new instance of Sort
func NewSort(id, value, description string) Sort {
	return Sort{ID: id, Value: value, Description: description}
}

// Sorting struct allows to sort a collection.
type Sorting struct {
	Sort      *Sort  `json:"sort,omitempty"`
	Available []Sort `json:"available,omitempty"`
}

// MarshalJSON supports json.Marshaler interface
func (v Sort) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1afec5d2EncodeGithubComIfreddyrondonCaptureAppListingSorting1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
