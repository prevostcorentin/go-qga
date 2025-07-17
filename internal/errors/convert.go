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

type ConvertError struct {
	wrappedError error
	kind         ConvertErrorKind
}

func NewConvertError(wrappedError error, kind ConvertErrorKind) *ConvertError {
	return &ConvertError{wrappedError: wrappedError, kind: kind}
}

func (_ *ConvertError) Domain() DomainType {
	return CodeGenerationDomain
}

func (convertError *ConvertError) Kind() string {
	return string(convertError.kind)
}

func (convertError *ConvertError) Unwrap() error {
	return convertError.wrappedError
}

func (convertError *ConvertError) Error() string {
	return formatErrorMessage(convertError)
}

type ConvertErrorKind string

const (
	UnknownPrimitiveType ConvertErrorKind = "Unknown primitive type"
	UnknownEntityType                     = "Unknown entity type"
)
