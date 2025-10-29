// Copyright 2025 The Archivista Contributors
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

package format

import (
	"context"
	"sync"
	"testing"
)

// mockHandler is a test implementation of Handler
type mockHandler struct {
	name       string
	detectFunc func(obj []byte) bool
}

func (m *mockHandler) Detect(obj []byte) bool {
	return m.detectFunc(obj)
}

func (m *mockHandler) Store(ctx context.Context, store Store, gitoid string, obj []byte) error {
	return nil
}

func TestRegisterHandler(t *testing.T) {
	// Clear handlers for isolated test
	handlersMu.Lock()
	handlers = nil
	handlersMu.Unlock()

	h1 := &mockHandler{name: "test1", detectFunc: func(obj []byte) bool { return false }}
	h2 := &mockHandler{name: "test2", detectFunc: func(obj []byte) bool { return false }}

	RegisterHandler(h1)
	RegisterHandler(h2)

	handlersMu.RLock()
	defer handlersMu.RUnlock()

	if len(handlers) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(handlers))
	}
}

func TestGetHandler_FirstMatchWins(t *testing.T) {
	// Clear handlers for isolated test
	handlersMu.Lock()
	handlers = nil
	handlersMu.Unlock()

	h1 := &mockHandler{
		name:       "test1",
		detectFunc: func(obj []byte) bool { return true }, // Always matches
	}
	h2 := &mockHandler{
		name:       "test2",
		detectFunc: func(obj []byte) bool { return true }, // Also matches
	}

	RegisterHandler(h1)
	RegisterHandler(h2)

	handler, ok := GetHandler([]byte("test"))
	if !ok {
		t.Fatal("expected to find a handler")
	}

	// Should get h1 because it was registered first
	if handler.(*mockHandler).name != "test1" {
		t.Errorf("expected h1 to be returned (first match wins), got %s", handler.(*mockHandler).name)
	}
}

func TestGetHandler_NoMatch(t *testing.T) {
	// Clear handlers for isolated test
	handlersMu.Lock()
	handlers = nil
	handlersMu.Unlock()

	h1 := &mockHandler{
		name:       "test1",
		detectFunc: func(obj []byte) bool { return false }, // Never matches
	}

	RegisterHandler(h1)

	handler, ok := GetHandler([]byte("test"))
	if ok {
		t.Error("expected no handler to be found")
	}
	if handler != nil {
		t.Error("expected nil handler when not found")
	}
}

func TestRegisterHandler_Concurrent(t *testing.T) {
	// Clear handlers for isolated test
	handlersMu.Lock()
	handlers = nil
	handlersMu.Unlock()

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrently register handlers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			h := &mockHandler{
				name:       "concurrent-handler",
				detectFunc: func(obj []byte) bool { return false },
			}
			RegisterHandler(h)
		}(i)
	}

	wg.Wait()

	handlersMu.RLock()
	count := len(handlers)
	handlersMu.RUnlock()

	if count != numGoroutines {
		t.Errorf("expected %d handlers after concurrent registration, got %d", numGoroutines, count)
	}
}

func TestGetHandler_Concurrent(t *testing.T) {
	// Clear handlers for isolated test
	handlersMu.Lock()
	handlers = nil
	handlersMu.Unlock()

	h1 := &mockHandler{
		name:       "test1",
		detectFunc: func(obj []byte) bool { return true },
	}
	RegisterHandler(h1)

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrently get handlers
	errors := make(chan error, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			handler, ok := GetHandler([]byte("test"))
			if !ok {
				errors <- nil
			}
			if handler == nil {
				errors <- nil
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("concurrent GetHandler failed: %v", err)
		}
	}
}
