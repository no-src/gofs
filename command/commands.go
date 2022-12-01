package command

import "github.com/no-src/log"

// Commands the commands that contain the init, actions and clear stages
type Commands struct {
	Name    string
	Init    []Command
	Actions []Command
	Clear   []Command
}

// Exec execute the commands in order of stages
func (cs *Commands) Exec() (err error) {
	if err = cs.ExecInit(); err != nil {
		return err
	}
	if err = cs.ExecActions(); err != nil {
		return err
	}
	return cs.ExecClear()
}

// ExecInit execute the init commands
func (cs *Commands) ExecInit() (err error) {
	return cs.exec("init", cs.Init)
}

// ExecActions execute the actions commands
func (cs *Commands) ExecActions() (err error) {
	return cs.exec("actions", cs.Actions)
}

// ExecClear execute the clear commands
func (cs *Commands) ExecClear() (err error) {
	return cs.exec("clear", cs.Clear)
}

func (cs *Commands) exec(stage string, commands []Command) (err error) {
	for i, c := range commands {
		if err = c.Exec(); err != nil {
			log.Error(err, "[%s] [%d] execute failed", stage, i+1)
			return err
		}
	}
	return nil
}
