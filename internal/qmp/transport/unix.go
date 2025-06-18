package transport

import (
	"bufio"
	"net"

	. "github.com/prevostcorentin/go-qga/internal/errors"
)

type unixTransport struct {
	path       string
	connection net.Conn
	pipe       *bufio.ReadWriter
}

func (transport *unixTransport) Connect() *TransportError {
	var err error
	if transport.connection, err = net.Dial("unix", transport.path); err != nil {
		return NewTransportError(err, Connect)
	}
	transport.pipe = bufio.NewReadWriter(
		bufio.NewReader(transport.connection),
		bufio.NewWriter(transport.connection),
	)
	return nil
}

func (transport *unixTransport) Write(bytes []byte) error {
	if _, err := transport.pipe.Write(bytes); err != nil {
		return NewTransportError(err, Write)
	}
	if err := transport.pipe.Writer.Flush(); err != nil {
		return NewTransportError(err, Flush)
	}
	return nil
}

func (transport *unixTransport) Read() ([]byte, error) {
	var err error
	var bytes []byte
	if bytes, err = transport.pipe.ReadBytes('\n'); err == nil {
		return bytes, nil
	}
	return nil, NewTransportError(err, Read)
}

func (transport *unixTransport) Path() string {
	return transport.path
}

func (transport *unixTransport) Close() error {
	if err := transport.pipe.Writer.Flush(); err != nil {
		return NewTransportError(err, Flush)
	}
	if err := transport.connection.Close(); err != nil {
		return NewTransportError(err, Close)
	}
	return nil
}
