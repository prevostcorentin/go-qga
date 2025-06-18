// Copyright 2025 PREVOST Corentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qmp

import (
	"encoding/json"
	"fmt"

	. "github.com/prevostcorentin/go-qga/internal/errors"
)

type commandExecutor struct {
	connection QmpConnection
}

func NewExecutor(connection QmpConnection) *commandExecutor {
	return &commandExecutor{connection: connection}
}

func (executor *commandExecutor) Run(command Command) (any, QgaError) {
	marshalled := struct {
		Execute   string `json:"execute"`
		Arguments any    `json:"arguments,omitempty"`
	}{
		Execute:   command.Execute(),
		Arguments: command.Arguments(),
	}
	marshalledBytes, marshalErr := json.Marshal(marshalled)
	if marshalErr != nil {
		return nil, NewCodecError(marshalErr, Marshal)
	}
	marshalledBytes = append(marshalledBytes, '\n') // Marshalling does not suffix objects with a line feed
	responseBytes, sendErr := executor.connection.Send(marshalledBytes)
	if sendErr != nil {
		return nil, sendErr
	}
	return executor.unmarshalCommandResponse(responseBytes, command)
}

func (executor *commandExecutor) unmarshalCommandResponse(bytes []byte, command Command) (any, QgaError) {
	typedResponse := command.Response()
	var root map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &root); err != nil {
		return nil, NewCodecError(err, Unmarshal)
	}
	if raw, ok := root["return"]; ok {
		if err := json.Unmarshal(raw, &typedResponse); err != nil {
			return nil, NewCodecError(err, Unmarshal)
		}
		return typedResponse, nil
	}
	return typedResponse, NewCodecError(fmt.Errorf(`missing "return" field in QGA response`), Key)
}
