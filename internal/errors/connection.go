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
