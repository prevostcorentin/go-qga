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

package qmp_test

import (
	"errors"
	"net"
	"os"
	"testing"

	qmpErrors "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp"
)

type fakeTransport struct{}

func (_ *fakeTransport) Connect() error {
	return nil
}

func (_ *fakeTransport) Close() error {
	return nil
}

func (_ *fakeTransport) Path() string {
	return "fakeTransport"
}

func (_ *fakeTransport) Read() ([]byte, error) {
	return nil, nil
}

func (_ *fakeTransport) Write(bytes []byte) error {
	return nil
}

func TestOpenUnexistingSocket(t *testing.T) {
	if _, err := qmp.Open("/that/path/will/not/exist/for/sure", &fakeTransport{}); err == nil {
		t.Error("socket does not exist. it should have raised an error but didn't.")
	}
}

type noWriteTransport struct{}

func (_ *noWriteTransport) Connect() error {
	return nil
}

func (_ *noWriteTransport) Close() error {
	return nil
}

func (_ *noWriteTransport) Path() string {
	return "fakeTransport"
}

func (_ *noWriteTransport) Read() ([]byte, error) {
	return nil, nil
}

func (_ *noWriteTransport) Write(bytes []byte) error {
	return errors.New("can't write anything. I am malfunctioning.")
}

func TestSendWriteMalfunction(t *testing.T) {
	socketPath, transport := buildSocketPath(), noWriteTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("couldn't open socket pipe %s", socketPath)
	}
	defer listener.Close()
	defer os.Remove(socketPath)
	socket, err := qmp.Open(socketPath, &transport)
	if err != nil {
		t.Fatal("the socket should open here")
	}
	_, err = socket.Send([]byte("no data"))
	if err == nil {
		t.Error("there should have been an error here")
	}
	if err.(*qmp.SocketError).Domain() != qmpErrors.SocketDomain {
		t.Errorf(`wrong error domain "%s". should have been "Socket"`, err.(*qmp.SocketError).Domain())
	}
	if err.(*qmp.SocketError).Kind() != qmp.SendErrorType {
		t.Errorf(`wrong error kind "%s". should have been "Send"`, err.(*qmp.SocketError).Kind())
	}
}

type noReadTransport struct{}

func (_ *noReadTransport) Connect() error {
	return nil
}

func (_ *noReadTransport) Close() error {
	return nil
}

func (_ *noReadTransport) Path() string {
	return "fakeTransport"
}

func (_ *noReadTransport) Read() ([]byte, error) {
	return nil, errors.New("can't read anything. I am malfunctioning.")
}

func (_ *noReadTransport) Write(bytes []byte) error {
	return nil
}

func TestSendReadMalfunction(t *testing.T) {
	socketPath, transport := buildSocketPath(), noWriteTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("couldn't open socket pipe %s", socketPath)
	}
	defer listener.Close()
	defer os.Remove(socketPath)
	socket, err := qmp.Open(socketPath, &transport)
	if err != nil {
		t.Fatal("the socket should open here")
	}
	_, err = socket.Send([]byte("no data"))
	if err == nil {
		t.Error("there should have been an error here")
	}
	if err.(*qmp.SocketError).Domain() != qmpErrors.SocketDomain {
		t.Errorf(`wrong error domain "%s". should have been "Socket"`, err.(*qmp.SocketError).Domain())
	}
	if err.(*qmp.SocketError).Kind() != qmp.SendErrorType {
		t.Errorf(`wrong error type "%s". should have been "Send"`, err.(*qmp.SocketError).Domain())
	}
}

type malfunctioningTransport struct{}

func (transport *malfunctioningTransport) Connect() error {
	return errors.New("malfunctioning transport")
}

func (transport *malfunctioningTransport) Close() error {
	return errors.New("malfunctioning transport")
}

func (transport *malfunctioningTransport) Path() string {
	return "malfunctioning transport"
}

func (transport *malfunctioningTransport) Read() ([]byte, error) {
	return nil, errors.New("malfunctioning transport")
}

func (transport *malfunctioningTransport) Write(_ []byte) error {
	return errors.New("malfunctioning transport")
}

func TestOpenMalfunctioningSocket(t *testing.T) {
	socketPath, transport := buildSocketPath(), malfunctioningTransport{}
	if _, err := qmp.Open(socketPath, &transport); err == nil {
		t.Error("socket should not open.")
	}
}
