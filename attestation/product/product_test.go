// Copyright 2021 The Witness Contributors
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

package product

import (
	"archive/tar"
	"bytes"
	"crypto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromDigestMap(t *testing.T) {
	testDigest, err := cryptoutil.CalculateDigestSetFromBytes([]byte("test"), []crypto.Hash{crypto.SHA256})
	assert.NoError(t, err)
	testDigestSet := make(map[string]cryptoutil.DigestSet)
	testDigestSet["test"] = testDigest
	result := fromDigestMap(testDigestSet)
	assert.Len(t, result, 1)
	digest := result["test"].Digest
	assert.True(t, digest.Equal(testDigest))
}

func TestAttestorName(t *testing.T) {
	a := New()
	assert.Equal(t, a.Name(), Name)
}

func TestAttestorType(t *testing.T) {
	a := New()
	assert.Equal(t, a.Type(), Type)
}

func TestAttestorRunType(t *testing.T) {
	a := New()
	assert.Equal(t, a.RunType(), RunType)
}

func TestAttestorAttest(t *testing.T) {
	a := New()
	testDigest, err := cryptoutil.CalculateDigestSetFromBytes([]byte("test"), []crypto.Hash{crypto.SHA256})
	if err != nil {
		t.Errorf("Failed to calculate digest set from bytes: %v", err)
	}

	testDigestSet := make(map[string]cryptoutil.DigestSet)
	testDigestSet["test"] = testDigest
	a.baseArtifacts = testDigestSet
	ctx, err := attestation.NewContext([]attestation.Attestor{a})
	require.NoError(t, err)
	require.NoError(t, a.Attest(ctx))
}

func TestGetFileContentType(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a temporary text file.
	textFile, err := os.CreateTemp(tempDir, "test-*.txt")
	require.NoError(t, err)
	defer os.Remove(textFile.Name())
	_, err = textFile.WriteString("This is a test file.")
	require.NoError(t, err)

	// Create a temporary PDF file with extension.
	pdfFile, err := os.CreateTemp(tempDir, "test-*")
	require.NoError(t, err)
	defer os.Remove(pdfFile.Name())

	//write to pdf so it has correct file signature 25 50 44 46 2D
	_, err = pdfFile.WriteAt([]byte{0x25, 0x50, 0x44, 0x46, 0x2D}, 0)

	require.NoError(t, err)

	// Create a temporary tar file with no extension.
	tarFile, err := os.CreateTemp(tempDir, "test-*")
	require.NoError(t, err)
	defer os.Remove(tarFile.Name())
	tarBuffer := new(bytes.Buffer)
	writer := tar.NewWriter(tarBuffer)
	header := &tar.Header{
		Name: "test.txt",
		Size: int64(len("This is a test file.")),
	}
	require.NoError(t, writer.WriteHeader(header))
	_, err = writer.Write([]byte("This is a test file."))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	_, err = tarFile.Write(tarBuffer.Bytes())
	require.NoError(t, err)

	// Open the temporary tar file using os.Open.
	tarFile, err = os.Open(tarFile.Name())
	require.NoError(t, err)

	// Define the test cases.
	tests := []struct {
		name     string
		file     *os.File
		expected string
	}{
		{"text file with extension", textFile, "text/plain; charset=utf-8"},
		{"PDF file with no extension", pdfFile, "application/pdf"},
		{"tar file with no extension", tarFile, "application/x-tar"},
	}

	// Run the test cases.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contentType, err := getFileContentType(test.file)
			require.NoError(t, err)
			require.Equal(t, test.expected, contentType)
		})
	}
}

func TestIncludeExcludeGlobs(t *testing.T) {
	workingDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(workingDir, "subdir"), 0777))
	files := []string{
		filepath.Join(workingDir, "test.txt"),
		filepath.Join(workingDir, "test.exe"),
		filepath.Join(workingDir, "subdir", "test.txt"),
		filepath.Join(workingDir, "subdir", "test.exe"),
	}

	for _, file := range files {
		f, err := os.Create(file)
		require.NoError(t, err)
		require.NoError(t, f.Close())
	}

	tests := []struct {
		name             string
		includeGlob      string
		excludeGlob      string
		expectedSubjects []string
	}{
		{"match all", "*", "", []string{"test.txt", "test.exe", filepath.Join("subdir", "test.txt"), filepath.Join("subdir", "test.exe")}},
		{"include only exes", "*.exe", "", []string{"test.exe", filepath.Join("subdir", "test.exe")}},
		{"exclude exes", "*", "*.exe", []string{"test.txt", filepath.Join("subdir", "test.txt")}},
		{"include only files in subdir", "subdir/*", "", []string{filepath.Join("subdir", "test.txt"), filepath.Join("subdir", "test.exe")}},
		{"exclude files in subdir", "*", "subdir/*", []string{"test.txt", "test.exe"}},
		{"include nothing", "", "", []string{}},
		{"exclude everything", "", "*", []string{}},
	}

	assertSubjsMatch := func(t *testing.T, subjects map[string]cryptoutil.DigestSet, expected []string) {
		subjectPaths := make([]string, 0, len(subjects))
		for path := range subjects {
			subjectPaths = append(subjectPaths, strings.TrimPrefix(path, "file:"))
		}

		assert.ElementsMatch(t, subjectPaths, expected)
	}

	t.Run("default include all", func(t *testing.T) {
		ctx, err := attestation.NewContext([]attestation.Attestor{}, attestation.WithWorkingDir(workingDir))
		require.NoError(t, err)
		a := New()
		require.NoError(t, a.Attest(ctx))
		assertSubjsMatch(t, a.Subjects(), []string{"test.txt", "test.exe", filepath.Join("subdir", "test.txt"), filepath.Join("subdir", "test.exe")})
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, err := attestation.NewContext([]attestation.Attestor{}, attestation.WithWorkingDir(workingDir))
			require.NoError(t, err)
			a := New()
			WithIncludeGlob(test.includeGlob)(a)
			WithExcludeGlob(test.excludeGlob)(a)
			require.NoError(t, a.Attest(ctx))
			assertSubjsMatch(t, a.Subjects(), test.expectedSubjects)
		})
	}
}
