package loader

import (
	"fmt"
)

// TaskConfig the config of tasks
type TaskConfig struct {
	// Tasks the task list
	Tasks []*TaskItem `json:"tasks" yaml:"tasks"`
}

// TaskItem a task item
type TaskItem struct {
	// Name a unique task name
	Name string `json:"name" yaml:"name"`
	// Conf the source of task config
	Conf string `json:"conf" yaml:"conf"`
	// Labels it can only acquire the current task if the client matches all the labels
	Labels []string `json:"labels" yaml:"labels"`
	// AllowIP the current task only allows the specified ip to access
	AllowIP []string `json:"allow_ip" yaml:"allow_ip"`
}

func (c *TaskConfig) verify() error {
	m := make(map[string]struct{})
	for _, t := range c.Tasks {
		if _, ok := m[t.Name]; ok {
			return fmt.Errorf("%w => %s", errDuplicateTask, t.Name)
		}
		m[t.Name] = struct{}{}
	}
	return nil
}
