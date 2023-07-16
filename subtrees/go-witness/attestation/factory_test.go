// Copyright 2023 The Archivist Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attestation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	attestors := []dummyAttestor{
		{
			name:          "prerun",
			predicateType: "https://witness.dev/test/prerun",
			runType:       PreMaterialRunType,
		}, {
			name:          "execute",
			predicateType: "https://witness.dev/test/execute",
			runType:       ExecuteRunType,
		},
		{
			name:          "post",
			predicateType: "https://witness.dev/test/post",
			runType:       PostProductRunType,
		},
	}

	for _, attestor := range attestors {
		RegisterAttestation(attestor.name, attestor.predicateType, attestor.runType, func() Attestor { return &attestor })
	}

	for _, attestor := range attestors {
		factory, ok := FactoryByType(attestor.predicateType)
		require.True(t, ok)
		otherFactory, ok := FactoryByName(attestor.name)
		require.True(t, ok)
		assert.Equal(t, factory(), otherFactory())
	}
}

type dummyAttestor struct {
	name          string
	predicateType string
	runType       RunType
	intOpt        int
	strOpt        string
	strSliceOpt   []string
}

func (a *dummyAttestor) Name() string {
	return a.name
}

func (a *dummyAttestor) Type() string {
	return a.predicateType
}

func (a *dummyAttestor) RunType() RunType {
	return a.runType
}

func (a *dummyAttestor) Attest(*AttestationContext) error {
	return nil
}
