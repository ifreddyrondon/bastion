package paging

import "github.com/mailru/easyjson/jwriter"

// Paging struct allows to do pagination into a collection.
type Paging struct {
	MaxAllowedLimit int   `json:"max_allowed_limit"`
	Limit           int   `json:"limit"`
	Offset          int64 `json:"offset"`
	Total           int64 `json:"total,omitempty"`
}

// MarshalJSON supports json.Marshaler interface
func (v Paging) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3c4140EncodeGithubComIfreddyrondonCaptureAppListingPaging(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
