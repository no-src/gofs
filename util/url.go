package util

import "net/url"

func ValuesEncode(k, v string) string {
	values := url.Values{}
	values.Add(k, v)
	return values.Encode()
}
