// Copyright 2024 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attestationcollection

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/in-toto/archivista/ent"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Define a mock parser function
	mockParser := func(ctx context.Context, tx *ent.Tx, attestation *ent.Attestation, attestationType string, message json.RawMessage) error {
		return nil
	}

	// Register the mock parser
	Register("mockType", mockParser)

	// Check if the parser is registered
	registeredParser, exists := registeredParsers["mockType"]
	var typedParser AttestationParser = mockParser
	assert.True(t, exists, "Parser should be registered")
	assert.IsType(t, typedParser, registeredParser, "Registered parser should match the mock parser")
}
