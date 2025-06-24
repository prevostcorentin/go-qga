package errors

type CollectErrorKind string

type CollectError struct {
	wrappedError error
	kind         CollectErrorKind
}

func NewCollectError(wrappedError error, errorKind CollectErrorKind) *CollectError {
	return &CollectError{wrappedError: wrappedError, kind: errorKind}
}

func (err *CollectError) Domain() DomainType {
	return CodeGenerationDomain
}

func (err *CollectError) Kind() string {
	return string(err.kind)
}

func (err *CollectError) Unwrap() error {
	return err.wrappedError
}

func (err *CollectError) Error() string {
	return formatErrorMessage(err)
}

const (
	MalformedSchema CollectErrorKind = "Malformed schema"
	Unknown                          = "Unknown type"
)
