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

type Socket interface {
	Connect(path string) *SocketError
	Send(bytes []byte) ([]byte, *SocketError)
	Close() error
}

type socket struct {
	transport transport.Transport
}

func Open(path string, transport transport.Transport) (Socket, *SocketError) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		errorReason := fmt.Errorf(`socket "%s" does not exist`, path)
		return nil, NewSocketError(errorReason, ConnectErrorType)
	}
	if err := transport.Connect(); err != nil {
		return nil, NewSocketError(err, ConnectErrorType)
	}
	socket := socket{transport: transport}
	if err := socket.Connect(path); err != nil {
		return nil, NewSocketError(err, ConnectErrorType)
	}
	return &socket, nil
}

func (socket *socket) Connect(path string) *SocketError {
	return socket.consumeBanner()
}

func (socket *socket) consumeBanner() *SocketError {
	// TODO: Use the banner to gather agent capabilities (could result in client code generation ?)
	if _, err := socket.transport.Read(); err != nil {
		return NewSocketError(err, ReadErrorType)
	}
	return nil
}

func (socket *socket) Send(bytes []byte) ([]byte, *SocketError) {
	if err := socket.transport.Write(bytes); err != nil {
		return nil, NewSocketError(err, SendErrorType)
	}
	bytes, err := socket.transport.Read()
	if err != nil {
		return bytes, NewSocketError(err, SendErrorType)
	}
	return bytes, nil
}

func (socket *socket) Close() error {
	if err := socket.transport.Close(); err != nil {
		return NewSocketError(err, CloseErrorType)
	}
	return nil
}
