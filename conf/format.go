package conf

import "strings"

// Format the supported config format
type Format struct {
	name       string
	extensions []string
}

// NewFormat create a new format info
func NewFormat(name string, extensions ...string) Format {
	return Format{
		name:       name,
		extensions: extensions,
	}
}

var (
	// JsonFormat the json format config
	JsonFormat = NewFormat("json", ".json")
	// YamlFormat the yaml format config
	YamlFormat = NewFormat("yaml", ".yaml", ".yml")
)

// MatchExt is the current extension matches the format
func (f Format) MatchExt(ext string) bool {
	ext = strings.ToLower(ext)
	for _, v := range f.extensions {
		if v == ext {
			return true
		}
	}
	return false
}

// Name return the format name
func (f Format) Name() string {
	return f.name
}
