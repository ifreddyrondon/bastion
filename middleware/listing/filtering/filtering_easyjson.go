package filtering

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

func easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering(out *jwriter.Writer, in Value) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	if in.Result != 0 {
		const prefix string = ",\"result\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Result))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Value) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

func easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering1(out *jwriter.Writer, in Filtering) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Filters) != 0 {
		const prefix string = ",\"filters\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v3, v4 := range in.Filters {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.Raw((v4).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	if len(in.Available) != 0 {
		const prefix string = ",\"available\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v5, v6 := range in.Available {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Raw((v6).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

func easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering2(out *jwriter.Writer, in Filter) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"values\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Values == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Values {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.Raw((v9).MarshalJSON())
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Filter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson8834d2f0EncodeGithubComIfreddyrondonCaptureAppListingFiltering2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}
