package command

// Config the config structure defined a series of commands
type Config struct {
	Name string `yaml:"name"`
	// IncludePath the root directory of the include files, include path is the directory of current config file by default
	IncludePath string   `yaml:"include_path"`
	Include     []string `yaml:"include"`
	Init        []Action `yaml:"init"`
	Actions     []Action `yaml:"actions"`
	Clear       []Action `yaml:"clear"`
}

// Action contain the command action name and some parameters that current command needed
type Action map[string]any
