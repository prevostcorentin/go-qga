// Copyright 2025 PREVOST Corentin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
