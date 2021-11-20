package util

import "net/url"

// ValuesEncode parse the k,v to url param
func ValuesEncode(k, v string) string {
	values := url.Values{}
	values.Add(k, v)
	return values.Encode()
}
