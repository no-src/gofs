package util

import (
	"fmt"
	"strconv"
)

// String parse the v to string
func String(v interface{}) string {
	switch v.(type) {
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
		data, err := Marshal(v)
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
