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

package transport_test

import (
	"bufio"
	"net"
	"testing"

	. "github.com/prevostcorentin/go-qga/internal/errors"
	"github.com/prevostcorentin/go-qga/internal/qmp/transport"
	. "github.com/prevostcorentin/go-qga/internal/testing"
)

func TestUnexistingSocketFailure(t *testing.T) {
	unexistingSocketPath := "/this/socket/does/not/exist/for/sure"
	unixTransport := transport.NewTransport("unix", unexistingSocketPath)
	if unixTransport.Path() != unexistingSocketPath {
		t.Fatalf(`wrong transport path "%v". expected "%s"`, unixTransport.Path(), unexistingSocketPath)
	}
	var transportError *TransportError
	if transportError = unixTransport.Connect(); transportError == nil {
		t.Fatal("there should have been an error here")
	}
	if transportError.Domain() != TransportDomain {
		t.Fatalf(`wrong error domain "%v". expected "%s"`, transportError.Domain(), TransportDomain)
	}
	if transportError.Kind() != string(Connect) {
		t.Fatalf(`wrong error kind "%v". expected "%s"`, transportError.Kind(), Connect)
	}
}

type echoAgent struct {
	listener net.Listener
	done     chan struct{}
	t        *testing.T
	path     string
}

func newEchoAgent(t *testing.T) *echoAgent {
	return &echoAgent{t: t, done: make(chan struct{}), path: BuildSocketPath(t)}
}

func (agent *echoAgent) Path() string {
	return agent.path
}

func (agent *echoAgent) Start() {
	var listenerError error
	agent.listener, listenerError = net.Listen("unix", agent.Path())
	if listenerError != nil {
		agent.t.Fatalf("can't listen on %s: %v", agent.Path(), listenerError)
	}
	go func() {
		for {
			connection, acceptError := agent.listener.Accept()
			if acceptError != nil {
				select {
				case <-agent.done:
					return
				default:
					agent.t.Fatalf("can't accept connection: %v", acceptError)
				}
			}
			go func() {
				defer connection.Close()
				connectionReader, connectionWriter := bufio.NewReader(connection), bufio.NewWriter(connection)
				bytes, connectionError := connectionReader.ReadBytes('\n')
				if connectionError != nil {
					agent.t.Fatalf("can't read: %v", connectionError)
				}
				if _, writeError := connectionWriter.Write(bytes); writeError != nil {
					agent.t.Fatalf("can't write: %v", writeError)
				}
				if flushError := connectionWriter.Flush(); flushError != nil {
					agent.t.Fatalf("can't flush: %v", flushError)
				}
			}()
		}
	}()
}

func (agent *echoAgent) Stop() {
	close(agent.done)
	if err := agent.listener.Close(); err != nil {
		agent.t.Fatalf("can't close listener: %v", err)
	}
}

func TestReadWrite(t *testing.T) {
	agent := newEchoAgent(t)
	agent.Start()
	unixTransport := transport.NewTransport("unix", agent.Path())
	if err := unixTransport.Connect(); err != nil {
		t.Fatalf("while connecting socket: %v", err)
	}
	expectedResponse := []byte("some string\n")
	if writeError := unixTransport.Write(expectedResponse); writeError != nil {
		t.Fatalf("while writing: %v", writeError)
	}
	response, readError := unixTransport.Read()
	if readError != nil {
		t.Fatalf("while reading: %v", readError)
	}
	if string(response) != string(expectedResponse) {
		t.Errorf(`wrong response "%v". expected "%s"`, string(response), string(expectedResponse))
	}
	agent.Stop()
}

type closeConnectionAgent struct {
	listener net.Listener
	t        *testing.T
	done     chan struct{}
	path     string
}

func newCloseConnectionAgent(t *testing.T) *closeConnectionAgent {
	return &closeConnectionAgent{t: t, done: make(chan struct{}), path: BuildSocketPath(t)}
}

func (agent *closeConnectionAgent) Path() string {
	return agent.path
}

func (agent *closeConnectionAgent) Start() {
	var listenerError error
	agent.listener, listenerError = net.Listen("unix", agent.Path())
	if listenerError != nil {
		agent.t.Fatalf("can't listen on %s: %v", agent.Path(), listenerError)
	}
	go func() {
		for {
			connection, acceptError := agent.listener.Accept()
			if acceptError != nil {
				select {
				case <-agent.done:
					return
				default:
					agent.t.Fatalf("can't accept connection: %v", acceptError)
				}
			}
			connection.Read([]byte{})
			connection.Close()
		}
	}()
}

func (agent *closeConnectionAgent) Stop() {
	close(agent.done)
	if err := agent.listener.Close(); err != nil {
		agent.t.Fatalf("can't close listener: %v", err)
	}
}

func TestNoWrite(t *testing.T) {
	agent := newCloseConnectionAgent(t)
	agent.Start()
	transport := transport.NewTransport("unix", agent.Path())
	if connectError := transport.Connect(); connectError != nil {
		t.Fatalf("while connecting: %v", connectError)
	}
	largePayload := make([]byte, 1<<20) // 1 MiB of zeros
	var writeError error
	if writeError = transport.Write(largePayload); writeError == nil {
		t.Fatal("there should have been an error here")
	}
	transportError := writeError.(*TransportError)
	if transportError.Domain() != TransportDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, transportError.Domain(), TransportDomain)
	}
	if transportError.Kind() != Write {
		t.Errorf(`wrong error kind "%v". expected "%s"`, transportError.Kind(), Write)
	}
}

func TestNoRead(t *testing.T) {
	agent := newCloseConnectionAgent(t)
	agent.Start()
	transport := transport.NewTransport("unix", agent.Path())
	if connectError := transport.Connect(); connectError != nil {
		t.Fatalf("while connecting: %v", connectError)
	}
	var readError error
	if _, readError = transport.Read(); readError == nil {
		t.Fatal("there should have been an error here")
	}
	transportError := readError.(*TransportError)
	if transportError.Domain() != TransportDomain {
		t.Errorf(`wrong error domain "%v". expected "%s"`, transportError.Domain(), TransportDomain)
	}
	if transportError.Kind() != Read {
		t.Errorf(`wrong error kind "%v". expected "%s"`, transportError.Kind(), Read)
	}
}
