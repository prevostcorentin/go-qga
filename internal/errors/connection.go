package errors

type QmpConnectionError struct {
	wrappedError error
	kind         QmpConnectionErrorKind
}

func NewQmpConnectionError(wrappedError error, errorType QmpConnectionErrorKind) *QmpConnectionError {
	return &QmpConnectionError{wrappedError: wrappedError, kind: errorType}
}

func (err *QmpConnectionError) Domain() DomainType {
	return QmpConnectionDomain
}

func (err *QmpConnectionError) Kind() string {
	return string(err.kind)
}

func (err *QmpConnectionError) Unwrap() error {
	return err.wrappedError
}

func (connectionError *QmpConnectionError) Error() string {
	return formatErrorMessage(connectionError)
}

type QmpConnectionErrorKind string

const (
	UnknownErrorKind QmpConnectionErrorKind = "Unknown"
	ConnectErrorKind                        = "Connect"
	SendErrorKind                           = "Send"
	ReadErrorKind                           = "Read"
	CloseErrorKind                          = "Close"
)
