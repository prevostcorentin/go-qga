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

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp"
	. "github.com/prevostcorentin/go-qga/internal/testing"
)

type fakeTransport struct{}

func (_ *fakeTransport) Connect() *TransportError {
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

func (_ *noWriteTransport) Connect() *TransportError {
	return nil
}

func (_ *noWriteTransport) Close() error {
	return nil
}

func (_ *noWriteTransport) Path() string {
	return "noWriteTransport"
}

func (_ *noWriteTransport) Read() ([]byte, error) {
	return nil, nil
}

func (_ *noWriteTransport) Write(bytes []byte) error {
	return errors.New("can't write anything. I am malfunctioning.")
}

func TestSendWriteMalfunction(t *testing.T) {
	socketPath, transport := BuildSocketPath(t), noWriteTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, listenErr := net.Listen("unix", socketPath)
	if listenErr != nil {
		t.Fatalf("couldn't open socket pipe %s: %v", socketPath, listenErr)
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
	if err.Domain() != QmpConnectionDomain {
		t.Errorf(`wrong error domain "%s". should have been "Socket"`, err.Domain())
	}
	if err.Kind() != SendErrorKind {
		t.Errorf(`wrong error kind "%s". should have been "Send"`, err.Kind())
	}
}

type noBannerTransport struct{}

func (_ *noBannerTransport) Connect() *TransportError {
	return nil
}

func (_ *noBannerTransport) Close() error {
	return nil
}

func (_ *noBannerTransport) Path() string {
	return "noBannerTransport"
}

func (_ *noBannerTransport) Read() ([]byte, error) {
	return nil, errors.New("can't read anything. I am malfunctioning")
}

func (_ *noBannerTransport) Write(bytes []byte) error {
	return nil
}

func TestMalfunctioningConnect(t *testing.T) {
	socketPath, transport := BuildSocketPath(t), noBannerTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, listenErr := net.Listen("unix", socketPath)
	if listenErr != nil {
		t.Fatalf("couldn't open socket pipe %s: %v", socketPath, listenErr)
	}
	defer listener.Close()
	defer os.Remove(socketPath)
	_, err := qmp.Open(socketPath, &transport)
	if err == nil {
		t.Fatal("the socket should not open here")
	}
	if err.Domain() != QmpConnectionDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), QmpConnectionDomain)
	}
	if err.Kind() != ConnectErrorKind {
		t.Errorf(`wrong error kind "%v". expected "%s"`, err.Domain(), QmpConnectionDomain)
	}
}

type noReadTransport struct {
	bannerRead bool
}

func (_ *noReadTransport) Connect() *TransportError {
	return nil
}

func (_ *noReadTransport) Close() error {
	return nil
}

func (_ *noReadTransport) Path() string {
	return "noReadTransport"
}

func (transport *noReadTransport) Read() ([]byte, error) {
	if transport.bannerRead {
		return nil, errors.New("can't read anything. I am malfunctioning.")
	}
	transport.bannerRead = true
	return nil, nil
}

func (_ *noReadTransport) Write(bytes []byte) error {
	return nil
}

func TestSendReadMalfunction(t *testing.T) {
	socketPath, transport := BuildSocketPath(t), noReadTransport{bannerRead: false}
	// We should create the socket pipe here unless the connection won't open
	listener, listenErr := net.Listen("unix", socketPath)
	if listenErr != nil {
		t.Fatalf("couldn't open socket pipe %s: %v", socketPath, listenErr)
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
	if err.Domain() != QmpConnectionDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), QmpConnectionDomain)
	}
	if err.Kind() != SendErrorKind {
		t.Errorf(`wrong error type "%v". expected "%s"`, err.Domain(), SendErrorKind)
	}
}

type noConnectTransport struct{}

func (_ *noConnectTransport) Connect() *TransportError {
	return NewTransportError(errors.New("i am malfunctioning"), Connect)
}

func (_ *noConnectTransport) Close() error {
	return nil
}

func (_ *noConnectTransport) Path() string {
	return "noConnectTransport"
}

func (_ *noConnectTransport) Read() ([]byte, error) {
	return nil, nil
}

func (_ *noConnectTransport) Write(bytes []byte) error {
	return nil
}

func TestConnectMalfunction(t *testing.T) {
	socketPath, transport := BuildSocketPath(t), noConnectTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, listenErr := net.Listen("unix", socketPath)
	if listenErr != nil {
		t.Fatalf("couldn't open socket pipe %s: %v", socketPath, listenErr)
	}
	defer listener.Close()
	defer os.Remove(socketPath)
	if _, err := qmp.Open(socketPath, &transport); err != nil {
		if err.Domain() != QmpConnectionDomain {
			t.Errorf(`wrong error domain "%v". expected "%s"`, err.Domain(), QmpConnectionDomain)
		}
		if err.Kind() != ConnectErrorKind {
			t.Errorf(`wrong error kind "%v". expected "%s"`, err.Kind(), SendErrorKind)
		}
	} else {
		t.Errorf("there should have been an error here")
	}
}

type notClosingTransport struct{}

func (_ *notClosingTransport) Connect() *TransportError {
	return nil
}

func (_ *notClosingTransport) Close() error {
	return errors.New("can't close. I am malfunctioning")
}

func (_ *notClosingTransport) Path() string {
	return "notClosingTransport"
}

func (_ *notClosingTransport) Read() ([]byte, error) {
	return nil, nil
}

func (_ *notClosingTransport) Write(bytes []byte) error {
	return nil
}

func TestClosingMalfunction(t *testing.T) {
	socketPath, transport := BuildSocketPath(t), notClosingTransport{}
	// We should create the socket pipe here unless the connection won't open
	listener, listenErr := net.Listen("unix", socketPath)
	if listenErr != nil {
		t.Fatalf("couldn't open socket pipe %s: %v", socketPath, listenErr)
	}
	defer listener.Close()
	defer os.Remove(socketPath)
	socket, openErr := qmp.Open(socketPath, &transport)
	if openErr != nil {
		t.Fatalf(`while opening socket: %v`, openErr)
	}
	var err error
	if err = socket.Close(); err == nil {
		t.Fatal("should have not been closed")
	}
	qmpConnectionError := err.(*QmpConnectionError)
	if qmpConnectionError.Domain() != QmpConnectionDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, qmpConnectionError.Domain(), QmpConnectionDomain)
	}
	if qmpConnectionError.Kind() != CloseErrorKind {
		t.Errorf(`wrong error kind "%v". expected "%s"`, qmpConnectionError.Kind(), CloseErrorKind)
	}
}

type malfunctioningTransport struct{}

func (transport *malfunctioningTransport) Connect() *TransportError {
	return NewTransportError(errors.New("malfunctioning transport"), Connect)
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
	socketPath, transport := BuildSocketPath(t), malfunctioningTransport{}
	if _, err := qmp.Open(socketPath, &transport); err == nil {
		t.Error("socket should not open.")
	}
}
