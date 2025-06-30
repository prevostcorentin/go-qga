package collector

type Command struct {
	CommandName    string         `json:"command"`
	CommandData    map[string]any `json:"data"`
	CommandReturns string         `json:"returns,omitempty"`
}

func (command *Command) Name() string {
	return command.CommandName
}

func (command *Command) Fields() map[string]any {
	return command.CommandData
}

func (_ *Command) Type() string {
	return string(CommandType)
}

func (command *Command) Returns() string {
	return command.CommandReturns
}
