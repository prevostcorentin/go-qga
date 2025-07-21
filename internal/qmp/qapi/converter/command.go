package converter

import "github.com/prevostcorentin/go-qga/internal/qmp/qapi/collector"

type Command struct {
	name        string
	arguments   []*Field
	returnsType string
}

func NewCommand(rawCommand *collector.Command) *Command {
	commandArguments, argumentsCounter := make([]*Field, len(rawCommand.Arguments)), 0
	for argumentName, argumentType := range rawCommand.Arguments {
		commandArguments[argumentsCounter] = NewField(argumentName, argumentType)
		argumentsCounter = argumentsCounter + 1
	}
	return &Command{name: rawCommand.Name(), arguments: commandArguments, returnsType: rawCommand.Returns}
}

func (command *Command) Name() string {
	return command.name
}

func (command *Command) Arguments() []*Field {
	return command.arguments
}

func (command *Command) ReturnsType() string {
	return command.returnsType
}
