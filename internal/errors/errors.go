package errors

type QgaError interface {
	Domain() DomainType
	Kind() string
	Unwrap() error
	Error() string
}

type DomainType string

const (
	TransportDomain DomainType = "Transport"
	SocketDomain               = "Socket"
	ProtocolDomain             = "Protocol"
)
