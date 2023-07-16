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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigOptions(t *testing.T) {
	entryName := "optionTest"
	defaultIntVal := 50
	defaultStrVal := "default string"
	defaultStrSliceVal := []string{"d", "e", "f"}
	defaultBoolVal := true
	testOpts := []Configurer{
		IntConfigOption("someint", "some int", defaultIntVal, func(te *testEntity, v int) (*testEntity, error) {
			te.intOpt = v
			return te, nil
		}),
		StringConfigOption("somestring", "some string", defaultStrVal, func(te *testEntity, v string) (*testEntity, error) {
			te.strOpt = v
			return te, nil
		}),
		StringSliceConfigOption("someslice", "some slice", defaultStrSliceVal, func(te *testEntity, v []string) (*testEntity, error) {
			te.strSliceOpt = v
			return te, nil
		}),
		BoolConfigOption("somebool", "some bool", defaultBoolVal, func(te *testEntity, v bool) (*testEntity, error) {
			te.boolOpt = v
			return te, nil
		}),
	}

	testReg := New[*testEntity]()

	testReg.Register(entryName, func() *testEntity { return &testEntity{} }, testOpts...)
	te, err := testReg.NewEntity(entryName)
	require.NoError(t, err)
	assert.Equal(t, te, &testEntity{
		intOpt:      defaultIntVal,
		strOpt:      defaultStrVal,
		strSliceOpt: defaultStrSliceVal,
		boolOpt:     defaultBoolVal,
	})

	intVal := 100
	strVal := "test string"
	strSliceVal := []string{"a", "b", "c"}
	boolVal := false
	opts, ok := testReg.Options(entryName)
	require.True(t, ok)
	for _, opt := range opts {
		switch o := opt.(type) {
		case *ConfigOption[*testEntity, int]:
			_, err = o.Setter()(te, intVal)
		case *ConfigOption[*testEntity, string]:
			_, err = o.Setter()(te, strVal)
		case *ConfigOption[*testEntity, []string]:
			_, err = o.Setter()(te, strSliceVal)
		case *ConfigOption[*testEntity, bool]:
			_, err = o.Setter()(te, boolVal)
		default:
			err = errors.New("unknown config option")
		}

		require.NoError(t, err)
	}

	assert.Equal(t, te, &testEntity{
		intOpt:      intVal,
		strOpt:      strVal,
		strSliceOpt: strSliceVal,
		boolOpt:     boolVal,
	})
}

type testEntity struct {
	intOpt      int
	strOpt      string
	strSliceOpt []string
	boolOpt     bool
}
