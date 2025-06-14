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
