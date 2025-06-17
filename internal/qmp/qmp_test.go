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
	"testing"

	"github.com/prevostcorentin/go-qga/internal/qmp"
	"github.com/prevostcorentin/go-qga/internal/qmp/transport"
)

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
	cleanTestFolder()
	socketPath := buildSocketPath()
	RunFakeQmpGuestAgent(t, socketPath)
	transport := transport.NewTransport(transport.Unix, socketPath)
	qgaSocket, err := qmp.Open(socketPath, transport)
	if err != nil {
		t.Fatalf("while opening socket: %v", err)
	}
	defer qgaSocket.Close()

	command := hostNameCommand{}
	executor := qmp.NewExecutor(qgaSocket)
	response, err := executor.Run(command)
	if err != nil {
		t.Fatalf("while running command: %v", err)
	}
	typedResponse := response.(*hostNameResponse)
	if typedResponse.Name != "fake-vm" {
		t.Errorf(`vm name differs (got "%s", expecting "fake-vm")`, typedResponse.Name)
	}
}
