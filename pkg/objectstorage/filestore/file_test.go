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
package filestore_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/in-toto/archivista/pkg/objectstorage/filestore"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT FileStoreSuite
type UTFileStoreSuite struct {
	suite.Suite
	tempDir string
	payload []byte
}

func TestUTFileStoreSuite(t *testing.T) {
	suite.Run(t, new(UTFileStoreSuite))
}

func (ut *UTFileStoreSuite) SetupTest() {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filestore_test")
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.tempDir = tempDir
	ut.payload = []byte("test payload")
}

func (ut *UTFileStoreSuite) TearDownTest() {
	os.RemoveAll(ut.tempDir)
}
func (ut *UTFileStoreSuite) Test_Get() {

	store, _, err := filestore.New(context.Background(), ut.tempDir, "")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// Define a test payload
	payload := []byte("test payload")

	// Store the payload
	err = store.Store(context.Background(), "test_gitoid", payload)
	if err != nil {
		ut.FailNow(err.Error())
	}

	// Attempt storing at malicious payload location
	err = store.Store(context.Background(), "../../test_gitoid", payload)
	if err != nil && err != filepath.ErrBadPattern {
		ut.FailNowf("Failed to detect bad path: %v", err.Error())
	}

	// Retrieve the payload
	reader, err := store.Get(context.Background(), "test_gitoid")
	if err != nil {
		ut.FailNowf("Failed to retrieve payload: %v", err.Error())
	}
	defer reader.Close()

	// Read the payload from the reader
	retrievedPayload, err := io.ReadAll(reader)
	if err != nil {
		ut.FailNowf("Failed to read payload: %v", err.Error())
	}

	// Compare the retrieved payload with the original payload
	ut.Equal(string(retrievedPayload), string(payload))

	// Attempt to retrieve non-local payload
	_, err = store.Get(context.Background(), "/etc/passwd")
	if err != nil && err != filepath.ErrBadPattern {
		ut.FailNowf("Failed to detect bad path: %v", err.Error())
	}

}
