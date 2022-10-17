// Copyright 2022 The Archivist Contributors
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
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
)

type Store struct {
	prefix string
}

func New(ctx context.Context, directory string, address string) (*Store, <-chan error, error) {
	errCh := make(chan error)
	go func() {
		server := handlers.CompressHandler(http.FileServer(http.Dir(directory)))
		log.Fatalln(http.ListenAndServe(address, server))
		<-ctx.Done()
		close(errCh)
	}()

	return &Store{
		prefix: directory,
	}, errCh, nil
}

func (s *Store) Get(ctx context.Context, gitoid string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.prefix, gitoid+".json"))
}

func (s *Store) Store(ctx context.Context, gitoid string, payload []byte) error {
	return os.WriteFile(filepath.Join(s.prefix, gitoid+".json"), payload, 0644)
}
