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

	"github.com/prevostcorentin/go-qga/internal/qmp/transport"
)

type Socket interface {
	Connect(path string) error
	Send(bytes []byte) ([]byte, error)
	Close() error
}

type SocketError struct {
	WrappedError error
	errType      SocketErrorType
}

func (err *SocketError) Domain() string {
	return "Socket"
}

func (err *SocketError) Kind() string {
	return string(err.errType)
}

func (err *SocketError) Unwrap() error {
	return err.WrappedError
}

func (err *SocketError) Error() string {
	message := fmt.Sprintf("Error: %s => %v", err.Domain(), err.WrappedError)
	return message
}

type SocketErrorType string

const (
	ConnectErrorType SocketErrorType = "Connect"
	SendErrorType                    = "Send"
	ReadErrorType                    = "Read"
	CloseErrorType                   = "Close"
)

type socket struct {
	transport transport.Transport
}

func Open(path string, transport transport.Transport) (Socket, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf(`socket "%s" does not exist`, path)
	}
	if err := transport.Connect(); err != nil {
		return nil, &SocketError{WrappedError: err, errType: ConnectErrorType}
	}
	socket := socket{transport: transport}
	if err := socket.Connect(path); err != nil {
		return nil, &SocketError{WrappedError: err, errType: ConnectErrorType}
	}
	return &socket, nil
}

func (socket *socket) Connect(path string) error {
	return socket.consumeBanner()
}

func (socket *socket) consumeBanner() error {
	// TODO: Use the banner to gather agent capabilities (could result in client code generation ?)
	if _, err := socket.transport.Read(); err != nil {
		return &SocketError{WrappedError: err, errType: ReadErrorType}
	}
	return nil
}

func (socket *socket) Send(bytes []byte) ([]byte, error) {
	if err := socket.transport.Write(bytes); err != nil {
		return nil, &SocketError{WrappedError: err, errType: SendErrorType}
	}
	bytes, err := socket.transport.Read()
	if err != nil {
		return bytes, &SocketError{WrappedError: err, errType: SendErrorType}
	}
	return bytes, nil
}

func (socket *socket) Close() error {
	if err := socket.transport.Close(); err != nil {
		return &SocketError{WrappedError: err, errType: CloseErrorType}
	}
	return nil
}
