package monitor

import "github.com/no-src/nsgo/hashutil"

// ToHashValueMessageList convert the hashutil.HashValues to a HashValue array
func ToHashValueMessageList(hvs hashutil.HashValues) []*HashValue {
	var list []*HashValue
	for _, hv := range hvs {
		list = append(list, &HashValue{
			Offset: hv.Offset,
			Hash:   hv.Hash,
		})
	}
	return list
}
