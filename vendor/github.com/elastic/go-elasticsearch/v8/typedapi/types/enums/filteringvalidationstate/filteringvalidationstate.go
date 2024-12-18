// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/2f823ff6fcaa7f3f0f9b990dc90512d8901e5d64

// Package filteringvalidationstate
package filteringvalidationstate

import "strings"

// https://github.com/elastic/elasticsearch-specification/blob/2f823ff6fcaa7f3f0f9b990dc90512d8901e5d64/specification/connector/_types/Connector.ts#L186-L190
type FilteringValidationState struct {
	Name string
}

var (
	Edited = FilteringValidationState{"edited"}

	Invalid = FilteringValidationState{"invalid"}

	Valid = FilteringValidationState{"valid"}
)

func (f FilteringValidationState) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}

func (f *FilteringValidationState) UnmarshalText(text []byte) error {
	switch strings.ReplaceAll(strings.ToLower(string(text)), "\"", "") {

	case "edited":
		*f = Edited
	case "invalid":
		*f = Invalid
	case "valid":
		*f = Valid
	default:
		*f = FilteringValidationState{string(text)}
	}

	return nil
}

func (f FilteringValidationState) String() string {
	return f.Name
}
