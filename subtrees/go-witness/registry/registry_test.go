// Copyright 2023 The Witness Contributors
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

package registry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllEntries(t *testing.T) {
	testReg := New[*testEntity]()
	expectedEntries := []Entry[*testEntity]{}
	for i := 0; i < 10; i++ {
		expectedEntries = append(expectedEntries, testReg.Register(fmt.Sprintf("entity-%v", i), func() *testEntity { return &testEntity{} }))
	}

	allEntries := testReg.AllEntries()
	// I'd prefer to use assert.ElementsMatch here, but because entries contain function pointers it will never work:
	// https://github.com/stretchr/testify/issues/1146
	assert.Len(t, allEntries, len(expectedEntries))
}

func TestNewEntity(t *testing.T) {
	testReg := New[*testEntity]()
	entityName := "test"
	testReg.Register(entityName, func() *testEntity { return &testEntity{} })
	te, err := testReg.NewEntity(entityName)
	require.NoError(t, err)
	assert.Equal(t, te, &testEntity{})
	te, err = testReg.NewEntity("this doesn't exist")
	assert.Error(t, err)
	assert.Nil(t, te)
}

func TestOptionPrefix(t *testing.T) {
	entryName := "test-entity"
	optName := "test-option"
	optPrefix := "prefix"
	testReg := New[*testEntity](WithOptionPrefix[*testEntity](optPrefix))
	testReg.Register(entryName, func() *testEntity { return &testEntity{} },
		StringConfigOption(optName, "prefix test", "", func(te *testEntity, v string) (*testEntity, error) { return te, nil }),
	)

	expectedOptName := fmt.Sprintf("%v-%v", optPrefix, optName)
	opts, ok := testReg.Options(entryName)
	require.True(t, ok)
	require.Len(t, opts, 1)
	assert.Equal(t, opts[0].Name(), expectedOptName)
}

func TestNoOptionPrefix(t *testing.T) {
	entryName := "test-entity"
	optName := "test-option"
	testReg := New[*testEntity]()
	testReg.Register(entryName, func() *testEntity { return &testEntity{} },
		StringConfigOption(optName, "prefix test", "", func(te *testEntity, v string) (*testEntity, error) { return te, nil }),
	)

	opts, ok := testReg.Options(entryName)
	require.True(t, ok)
	require.Len(t, opts, 1)
	assert.Equal(t, opts[0].Name(), optName)
}
