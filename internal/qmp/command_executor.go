package qmp

import (
	"encoding/json"
	"fmt"
)

type commandExecutor struct {
	socket *Socket
}

func NewExecutor(socket *Socket) *commandExecutor {
	return &commandExecutor{socket: socket}
}

func (executor *commandExecutor) Run(command Command) (any, error) {
	marshalled := struct {
		Execute   string `json:"execute"`
		Arguments any    `json:"arguments,omitempty"`
	}{
		Execute:   command.Execute(),
		Arguments: command.Arguments(),
	}
	marshalledBytes, err := json.Marshal(marshalled)
	if err != nil {
		return nil, err
	}
	marshalledBytes = append(marshalledBytes, 0x0A) // Marshalling does not suffix objects with a line feed
	responseBytes, err := executor.socket.send(marshalledBytes)
	if err != nil {
		return nil, err
	}
	return executor.unmarshalCommandResponse(responseBytes, command)
}

func (executor *commandExecutor) unmarshalCommandResponse(bytes []byte, command Command) (any, error) {
	typedResponse := command.Response()
	var root map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &root); err != nil {
		return nil, err
	}
	if raw, ok := root["return"]; ok {
		err := json.Unmarshal(raw, &typedResponse)
		return typedResponse, err
	}
	return typedResponse, fmt.Errorf(`missing "return" field in QGA response`)
}
