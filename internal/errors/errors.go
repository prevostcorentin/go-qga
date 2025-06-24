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

import "fmt"

type QgaError interface {
	error
	Domain() DomainType
	Kind() string
	Unwrap() error
}

type DomainType string

const (
	TransportDomain      DomainType = "Transport"
	QmpConnectionDomain             = "Connection"
	ProtocolDomain                  = "Protocol"
	CodecDomain                     = "Codec"
	CodeGenerationDomain            = "Codegen"
)

func formatErrorMessage(err QgaError) string {
	message := fmt.Sprintf("Error: %s => %v", err.Domain(), err.Unwrap())
	return message
}
