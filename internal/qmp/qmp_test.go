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
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/prevostcorentin/go-qga/internal/qmp"
	"github.com/prevostcorentin/go-qga/internal/qmp/transport"
	. "github.com/prevostcorentin/go-qga/internal/testing"
)

type QmpBannerResponse struct {
	Qmp struct {
		Version struct {
			Qemu struct {
				Major string `json:"major"`
				Minor string `json:"minor"`
				Micro string `json:"micro"`
			} `json:"qemu"`
			Package string `json:"package"`
		} `json:"version"`
		Capabilities []any `json:"capabilities"`
	} `json:"QMP"`
}

type QmpHostnameResponse struct {
	Return struct {
		Name string `json:"name"`
	} `json:"return"`
}

type QmpCommand struct {
	Execute string `json:"execute"`
}

type QmpError struct {
	Error struct {
		Class       string `json:"class"`
		Description string `json:"desc"`
	}
}

type fakeGuestAgent struct {
	listener net.Listener
	done     chan struct{}
	t        *testing.T
	path     string
}

func newFakeGuestAgent(t *testing.T) *fakeGuestAgent {
	return &fakeGuestAgent{t: t, done: make(chan struct{}), path: BuildSocketPath(t)}
}

func (agent *fakeGuestAgent) Start() {
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
			go handleConnection(agent.t, connection)
		}
	}()
}

func handleConnection(t *testing.T, connection net.Conn) {
	defer connection.Close()
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	banner := QmpBannerResponse{}
	bytes, _ := json.Marshal(banner)
	fmt.Fprintln(writer, string(bytes))
	writer.Flush()

	line, err := reader.ReadBytes(0x0A)
	t.Logf("%d bytes received", len(line))
	if err != nil {
		return
	}
	var command QmpCommand
	if err := json.Unmarshal(line, &command); err != nil {
		t.Fatalf("unmarshalling command: %v", err)
	}
	var response any
	if command.Execute == "guest-get-host-name" {
		qmpResponse := &QmpHostnameResponse{}
		qmpResponse.Return.Name = "fake-vm"
		response = qmpResponse
	} else {
		response = QmpError{}
	}
	bytes, err = json.Marshal(response)
	if err != nil {
		t.Fatalf("marshalling response: %v", err)
	}
	fmt.Fprintln(writer, string(bytes))
	t.Logf("%d bytes sent", len(bytes))
	writer.Flush()
}

func (agent *fakeGuestAgent) Path() string {
	return agent.path
}

func (agent *fakeGuestAgent) Stop() {
	close(agent.done)
	if err := agent.listener.Close(); err != nil {
		agent.t.Fatalf("can't close listener: %v", err)
	}
}

type hostNameCommand struct{}

func (command hostNameCommand) Execute() string {
	return "guest-get-host-name"
}

func (command hostNameCommand) Arguments() any {
	return nil
}

func (command hostNameCommand) Response() any {
	return &hostNameResponse{}
}

type hostNameResponse struct {
	Name string
}

func TestHostnameCommand(t *testing.T) {
	agent := newFakeGuestAgent(t)
	socketPath := agent.Path()
	transport := transport.NewTransport(transport.Unix, socketPath)
	agent.Start()
	qgaSocket, openErr := qmp.Open(socketPath, transport)
	if openErr != nil {
		t.Fatalf("while opening socket: %v", openErr)
	}
	defer qgaSocket.Close()

	command := hostNameCommand{}
	executor := qmp.NewExecutor(qgaSocket)
	response, err := executor.Run(command)
	if err != nil {
		t.Fatalf("while running command: %v", err)
	}
	agent.Stop()
	typedResponse := response.(*hostNameResponse)
	if typedResponse.Name != "fake-vm" {
		t.Errorf(`vm name differs (got "%s", expecting "fake-vm")`, typedResponse.Name)
	}
}
