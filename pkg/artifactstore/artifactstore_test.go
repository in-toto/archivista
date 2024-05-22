// Copyright 2024 The Archivista Contributors
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

package artifactstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestConfigFile(t *testing.T, workingDir, version, distroFilePath, distroDigest string) string {
	testConfig := `artifacts:
  witness:
    versions:
      ` + version + `:
        description: test
        distributions:
          linux:
            filelocation: ` + distroFilePath + `
            sha256digest: ` + distroDigest
	testConfigFilePath := filepath.Join(workingDir, "config.yaml")
	testConfigFile, err := os.Create(testConfigFilePath)
	require.NoError(t, err)
	_, err = testConfigFile.WriteString(testConfig)
	require.NoError(t, err)
	require.NoError(t, testConfigFile.Close())
	return testConfigFilePath
}

func TestStore(t *testing.T) {
	workingDir := t.TempDir()
	testDistroFilePath := filepath.Join(workingDir, "test")
	testDistroDigest := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	testVersion := "v0.1.0"
	testDistroFile, err := os.Create(testDistroFilePath)
	require.NoError(t, err)
	_, err = testDistroFile.Write([]byte("test"))
	require.NoError(t, err)
	require.NoError(t, testDistroFile.Close())

	t.Run("all good", func(t *testing.T) {
		testArtifactName := "witness"
		testDistroName := "linux"
		testConfigFilePath := createTestConfigFile(t, workingDir, testVersion, testDistroFilePath, testDistroDigest)
		as, err := New(WithConfigFile(testConfigFilePath))
		require.NoError(t, err)
		artifacts := as.Artifacts()
		assert.Len(t, artifacts, 1)
		versions, ok := as.Versions(testArtifactName)
		assert.True(t, ok)
		assert.Len(t, versions, 1)
		version, ok := as.Version(testArtifactName, testVersion)
		assert.True(t, ok)
		assert.Len(t, version.Distributions, 1)
		testDistro, ok := as.Distribution(testArtifactName, testVersion, testDistroName)
		assert.True(t, ok)
		assert.Equal(t, testDistro.FileLocation, testDistroFilePath)
		assert.Equal(t, testDistro.SHA256Digest, testDistroDigest)
	})

	t.Run("wrong file path", func(t *testing.T) {
		testConfigFilePath := createTestConfigFile(t, workingDir, testVersion, "garbage", testDistroDigest)
		_, err := New(WithConfigFile(testConfigFilePath))
		assert.Error(t, err)

	})

	t.Run("wrong file digest", func(t *testing.T) {
		testConfigFilePath := createTestConfigFile(t, workingDir, testVersion, testDistroFilePath, "garbage")
		_, err := New(WithConfigFile(testConfigFilePath))
		assert.Error(t, err)
	})
}
