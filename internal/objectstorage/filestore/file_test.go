// Copyright 2022 The Archivista Contributors
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

	filestore "github.com/in-toto/archivista/internal/objectstorage/filestore"
)

func TestStore_Get(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filestore_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new file store
	store, _, err := filestore.New(context.Background(), tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create file store: %v", err)
	}

	// Define a test payload
	payload := []byte("test payload")

	// Store the payload
	err = store.Store(context.Background(), "test_gitoid", payload)
	if err != nil {
		t.Fatalf("Failed to store payload: %v", err)
	}

	// Attempt storing at malicious payload location
	err = store.Store(context.Background(), "../../test_gitoid", payload)
	if err != nil && err != filepath.ErrBadPattern {
		t.Errorf("Failed to detect bad path: %v", err)
	}

	// Retrieve the payload
	reader, err := store.Get(context.Background(), "test_gitoid")
	if err != nil {
		t.Errorf("Failed to retrieve payload: %v", err)
	}
	defer reader.Close()

	// Read the payload from the reader
	retrievedPayload, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read payload: %v", err)
	}

	// Compare the retrieved payload with the original payload
	if string(retrievedPayload) != string(payload) {
		t.Errorf("Retrieved payload does not match original payload")
	}

	// Attempt to retrieve non-local payload
	_, err = store.Get(context.Background(), "/etc/passwd")
	if err != nil && err != filepath.ErrBadPattern {
		t.Errorf("Failed to detect bad path: %v", err)
	}

}
