package command

// Config the config structure defined a series of commands
type Config struct {
	Name    string   `yaml:"name"`
	Actions []Action `yaml:"actions"`
}

// Action contain the command action name and some parameters that current command needed
type Action map[string]any
