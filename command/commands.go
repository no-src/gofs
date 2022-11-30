package command

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
	return cs.exec(cs.Init)
}

// ExecActions execute the actions commands
func (cs *Commands) ExecActions() (err error) {
	return cs.exec(cs.Actions)
}

// ExecClear execute the clear commands
func (cs *Commands) ExecClear() (err error) {
	return cs.exec(cs.Clear)
}

func (cs *Commands) exec(commands []Command) (err error) {
	for _, c := range commands {
		if err = c.Exec(); err != nil {
			return err
		}
	}
	return nil
}
