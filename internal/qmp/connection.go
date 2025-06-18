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
	"errors"
	"fmt"
	"os"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/transport"
)

type QmpConnection interface {
	Connect(path string) *QmpConnectionError
	Send(bytes []byte) ([]byte, *QmpConnectionError)
	Close() error
}

type connection struct {
	transport transport.Transport
}

func Open(path string, transport transport.Transport) (QmpConnection, *QmpConnectionError) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		errorReason := fmt.Errorf(`socket "%s" does not exist`, path)
		return nil, NewQmpConnectionError(errorReason, ConnectErrorKind)
	}
	if err := transport.Connect(); err != nil {
		return nil, NewQmpConnectionError(err, ConnectErrorKind)
	}
	socket := connection{transport: transport}
	if err := socket.Connect(path); err != nil {
		return nil, NewQmpConnectionError(err, ConnectErrorKind)
	}
	return &socket, nil
}

func (connection *connection) Connect(path string) *QmpConnectionError {
	return connection.consumeBanner()
}

func (connection *connection) consumeBanner() *QmpConnectionError {
	// TODO: Use the banner to gather agent capabilities (could result in client code generation ?)
	if _, err := connection.transport.Read(); err != nil {
		return NewQmpConnectionError(err, ReadErrorKind)
	}
	return nil
}

func (connection *connection) Send(bytes []byte) ([]byte, *QmpConnectionError) {
	if err := connection.transport.Write(bytes); err != nil {
		return nil, NewQmpConnectionError(err, SendErrorKind)
	}
	bytes, err := connection.transport.Read()
	if err != nil {
		return bytes, NewQmpConnectionError(err, SendErrorKind)
	}
	return bytes, nil
}

func (connection *connection) Close() error {
	if err := connection.transport.Close(); err != nil {
		return NewQmpConnectionError(err, CloseErrorKind)
	}
	return nil
}
