package errors

type TransportError struct {
	wrappedError error
	kind         TransportErrorKind
}

func NewTransportError(wrappedError error, kind TransportErrorKind) *TransportError {
	return &TransportError{wrappedError: wrappedError, kind: kind}
}

func (_ *TransportError) Domain() DomainType {
	return TransportDomain
}

func (transportError *TransportError) Kind() string {
	return string(transportError.kind)
}

func (transportError *TransportError) Unwrap() error {
	return transportError.wrappedError
}

func (err *TransportError) Error() string {
	return formatErrorMessage(err)
}

type TransportErrorKind string

const (
	Connect      TransportErrorKind = "Connect"
	Write                           = "Write"
	Read                            = "Read"
	Close                           = "Close"
	Flush                           = "Flush"
	Timeout                         = "Timeout"
	NotConnected                    = "Not connected"
)
