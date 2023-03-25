package yamlutil

import (
	"gopkg.in/yaml.v3"
)

var (
	// Marshal serializes the value provided into a YAML document.
	Marshal = marshal

	// Unmarshal decodes the first document found within the in byte slice
	// and assigns decoded values into the out value.
	Unmarshal = unmarshal
)

func marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
