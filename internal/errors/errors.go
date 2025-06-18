package errors

import "fmt"

type QgaError interface {
	error
	Domain() DomainType
	Kind() string
	Unwrap() error
}

type DomainType string

const (
	TransportDomain     DomainType = "Transport"
	QmpConnectionDomain            = "Connection"
	ProtocolDomain                 = "Protocol"
	CodecDomain                    = "Codec"
)

func formatErrorMessage(err QgaError) string {
	message := fmt.Sprintf("Error: %s => %v", err.Domain(), err.Unwrap())
	return message
}
