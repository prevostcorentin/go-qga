package qmp

import (
	"encoding/json"
)

type Command[T any, R any] struct {
	Execute   string `json:"execute"`
	Arguments *T     `json:"arguments,omitempty"`
}

func (command *Command[T, R]) Run(socket *Socket) (*R, error) {
	rawResponse, err := socket.send(command)
	if err != nil {
		return nil, err
	}
	var response R
	err = json.Unmarshal(rawResponse, &response)
	return &response, err
}
