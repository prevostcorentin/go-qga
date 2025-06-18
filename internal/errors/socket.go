package errors

import "fmt"

type SocketError struct {
	wrappedError error
	errorType    SocketErrorType
}

func NewSocketError(wrappedError error, errorType SocketErrorType) *SocketError {
	return &SocketError{wrappedError: wrappedError, errorType: errorType}
}

func (err *SocketError) Domain() string {
	return "Socket"
}

func (err *SocketError) Kind() string {
	return string(err.errorType)
}

func (err *SocketError) Unwrap() error {
	return err.wrappedError
}

func (err *SocketError) Error() string {
	message := fmt.Sprintf("Error: %s => %v", err.Domain(), err.wrappedError)
	return message
}

type SocketErrorType string

const (
	UnknownErrorType SocketErrorType = "Unknown"
	ConnectErrorType                 = "Connect"
	SendErrorType                    = "Send"
	ReadErrorType                    = "Read"
	CloseErrorType                   = "Close"
)
