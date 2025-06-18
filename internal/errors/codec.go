package errors

type CodecError struct {
	wrappedError error
	kind         CodecErrorKind
}

func NewCodecError(wrappedError error, kind CodecErrorKind) *CodecError {
	return &CodecError{wrappedError: wrappedError, kind: kind}
}

func (_ *CodecError) Domain() DomainType {
	return CodecDomain
}

func (codecError *CodecError) Kind() string {
	return string(codecError.kind)
}

func (codecError *CodecError) Unwrap() error {
	return codecError.wrappedError
}

func (codecError *CodecError) Error() string {
	return formatErrorMessage(codecError)
}

type CodecErrorKind string

const (
	Marshal   CodecErrorKind = "Marshal"
	Unmarshal CodecErrorKind = "Unmarshal"
	Type      CodecErrorKind = "Type"
	Key       CodecErrorKind = "Key"
)
