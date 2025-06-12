package qmp

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func (socket *Socket) send(object interface{}) ([]byte, error) {
	jsonQuery, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	if _, err = fmt.Fprintf(socket.pipe.Writer, "%s\n", jsonQuery); err != nil {
		return make([]byte, 0, 0), err
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
