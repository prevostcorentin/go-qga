package transport

import (
	. "github.com/prevostcorentin/go-qga/internal/errors"
)

type TransportType string

const (
	Unix TransportType = "unix"
)

type Transport interface {
	Connect() *TransportError
	Close() error
	Path() string
	Read() ([]byte, error)
	Write(bytes []byte) error
}

func NewTransport(transportType TransportType, path string) Transport {
	var transport Transport
	switch transportType {
	case Unix:
		transport = &unixTransport{path: path}
	}
	return transport
}
