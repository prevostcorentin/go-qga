package collector

type Command struct {
	CommandName string            `json:"command"`
	Arguments   map[string]string `json:"data"`
	Returns     string            `json:"returns,omitempty"`
}

func (command *Command) Name() string {
	return command.CommandName
}
