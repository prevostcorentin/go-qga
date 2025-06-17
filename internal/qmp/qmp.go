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
	"bufio"
	"net"
)

type Socket struct {
	connection net.Conn
	pipe       *bufio.ReadWriter
}

func (socket *Socket) Connect(path string) error {
	var err error
	socket.connection, err = net.Dial("unix", path)
	socket.pipe = bufio.NewReadWriter(
		bufio.NewReader(socket.connection),
		bufio.NewWriter(socket.connection),
	)
	socket.consumeBanner()
	return err
}

func (socket *Socket) consumeBanner() error {
	// TODO: Use the banner to gather agent capabilities (could result in client code generation ?)
	_, err := socket.pipe.ReadBytes(0x0A) // Just consume, do nothing with it
	return err
}

func (socket *Socket) send(bytes []byte) ([]byte, error) {
	if _, err := socket.pipe.Write(bytes); err != nil {
		return nil, err
	}
	socket.pipe.Flush()
	bytes, err := socket.pipe.ReadBytes(0x0A)
	return bytes, err
}

func (socket *Socket) Close() error {
	if err := socket.pipe.Flush(); err != nil {
		return err
	}
	if err := socket.connection.Close(); err != nil {
		return err
	}
	return nil
}
