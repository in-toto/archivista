// Copyright 2022 The Archivista Contributors
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

package filestore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
)

type Store struct {
	prefix string
}

// New creates a Store backed by directory. If address is not empty, the
// contents of directory are also served over HTTP on that address until ctx
// is done. The returned channel is closed once the store has shut down.
func New(ctx context.Context, directory string, address string) (*Store, <-chan error, error) {
	errCh := make(chan error)
	store := &Store{
		prefix: directory,
	}

	if address == "" {
		go func() {
			<-ctx.Done()
			close(errCh)
		}()

		return store, errCh, nil
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("could not listen on %s: %w", address, err)
	}

	server := &http.Server{
		Handler:      handlers.CompressHandler(http.FileServer(http.Dir(directory))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not serve on %s: %v\n", address, err)
		}
	}()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("File server shutdown failed: %+v", err)
		}

		close(errCh)
	}()

	return store, errCh, nil
}

func (s *Store) Get(ctx context.Context, gitoid string) (io.ReadCloser, error) {
	if filepath.IsLocal(gitoid) {
		return os.Open(filepath.Join(s.prefix, gitoid+".json"))
	} else {
		return nil, filepath.ErrBadPattern
	}
}

func (s *Store) Store(ctx context.Context, gitoid string, payload []byte) error {
	if filepath.IsLocal(gitoid) {
		return os.WriteFile(filepath.Join(s.prefix, gitoid+".json"), payload, 0o600)
	} else {
		return filepath.ErrBadPattern
	}
}
