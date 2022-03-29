package jsonutil

import "encoding/json"

// Marshal returns the JSON encoding of v
var Marshal = marshal

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an InvalidUnmarshalError.
var Unmarshal = unmarshal

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
