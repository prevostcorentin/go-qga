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
