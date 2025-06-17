package transport

type TransportType string

const (
	Unix TransportType = "unix"
)

type Transport interface {
	Connect() error
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

type TransportErrorType string

const (
	UnknownError TransportErrorType = "Unknown"
	ReadError                       = "Reading"
	WriteError                      = "Writing"
	ConnectError                    = "Connecting"
	FlushError                      = "Flushing"
	TimeoutError                    = "Timeout"
	CloseError                      = "Closing"
)

type TransportError interface {
	Type() TransportErrorType
}
