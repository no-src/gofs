package stringutil

import (
	"fmt"
	"github.com/no-src/gofs/util/jsonutil"
	"strconv"
	"strings"
)

// ToString support to convert the current instance to string value
type ToString interface {
	// String return current format info
	String() string
}

// String parse the v to string
func String(v interface{}) string {
	switch v.(type) {
	case ToString:
		return v.(ToString).String()
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case uint64:
		return strconv.FormatUint(v.(uint64), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case bool:
		return strconv.FormatBool(v.(bool))
	case error:
		return v.(error).Error()
	default:
		data, err := jsonutil.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(data)
	}
}

// Int64 parse the string to int64
func Int64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

// IsEmpty is empty or whitespace string
func IsEmpty(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}
