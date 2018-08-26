package listing

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonDe046902EncodeGithubComIfreddyrondonCaptureAppListing(out *jwriter.Writer, in Listing) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"paging\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Paging).MarshalJSON())
	}
	if in.Sorting != nil {
		const prefix string = ",\"sorting\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.Sorting).MarshalJSON())
	}
	if in.Filtering != nil {
		const prefix string = ",\"filtering\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((*in.Filtering).MarshalJSON())
	}
	out.RawByte('}')
}
