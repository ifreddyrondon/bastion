package listing

import (
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/mailru/easyjson/jwriter"
)

// Listing holds the info to perform filtering, sorting and paging over a collection.
type Listing struct {
	Paging    paging.Paging        `json:"paging,omitempty"`
	Sorting   *sorting.Sorting     `json:"sorting,omitempty"`
	Filtering *filtering.Filtering `json:"filtering,omitempty"`
}

// MarshalJSON supports json.Marshaler interface
func (v Listing) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonDe046902EncodeGithubComIfreddyrondonCaptureAppListing(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
