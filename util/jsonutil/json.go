package jsonutil

import "encoding/json"

var (
	// Marshal returns the JSON encoding of v
	Marshal = marshal

	// MarshalIndent is like Marshal but applies Indent to format the output
	MarshalIndent = marshalIndent

	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v. If v is nil or not a pointer,
	// Unmarshal returns an InvalidUnmarshalError.
	Unmarshal = unmarshal
)

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func marshalIndent(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

func unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
