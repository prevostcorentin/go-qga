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

type ResolveError struct {
	wrappedError error
	kind         ResolveErrorKind
}

func NewResolveError(wrappedError error, kind ResolveErrorKind) *ResolveError {
	return &ResolveError{wrappedError: wrappedError, kind: kind}
}

func (_ *ResolveError) Domain() DomainType {
	return CodeGenerationDomain
}

func (ResolveError *ResolveError) Kind() string {
	return string(ResolveError.kind)
}

func (ResolveError *ResolveError) Unwrap() error {
	return ResolveError.wrappedError
}

func (ResolveError *ResolveError) Error() string {
	return formatErrorMessage(ResolveError)
}

type ResolveErrorKind string

const (
	UnknownReference ResolveErrorKind = "Reference not found"
)
