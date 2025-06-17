package transport

import (
	"bufio"
	"fmt"
	"net"
)

type unixTransport struct {
	path       string
	connection net.Conn
	pipe       *bufio.ReadWriter
}

func (transport *unixTransport) Connect() error {
	var err error
	if transport.connection, err = net.Dial("unix", transport.path); err != nil {
		return &UnixTransportError{WrappedError: err, errType: ConnectError}
	}
	transport.pipe = bufio.NewReadWriter(
		bufio.NewReader(transport.connection),
		bufio.NewWriter(transport.connection),
	)
	return nil
}

func (transport *unixTransport) Write(bytes []byte) error {
	if _, err := transport.pipe.Write(bytes); err != nil {
		return &UnixTransportError{WrappedError: err, errType: WriteError}
	}
	return transport.pipe.Flush()
}

func (transport *unixTransport) Read() ([]byte, error) {
	var err error
	var bytes []byte
	if bytes, err = transport.pipe.ReadBytes(0x0A); err == nil {
		return bytes, nil
	}
	return nil, &UnixTransportError{WrappedError: err, errType: ReadError}
}

func (transport *unixTransport) Path() string {
	return transport.path
}

func (transport *unixTransport) Close() error {
	if err := transport.pipe.Flush(); err != nil {
		return &UnixTransportError{WrappedError: err, errType: FlushError}
	}
	if err := transport.connection.Close(); err != nil {
		return &UnixTransportError{WrappedError: err, errType: CloseError}
	}
	return nil
}

type UnixTransportError struct {
	WrappedError error
	errType      TransportErrorType
}

func (err *UnixTransportError) Type() TransportErrorType {
	return err.errType
}

func (err *UnixTransportError) Unwrap() error {
	return err.WrappedError
}

func (err *UnixTransportError) Error() string {
	message := fmt.Sprintf("%s error: %v\n", err.errType, err.WrappedError)
	return message
}
