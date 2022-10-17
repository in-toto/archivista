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

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/edwarnicke/gitoid"
	"github.com/gorilla/mux"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	archivistapi "github.com/testifysec/archivist-api"
)

type Server struct {
	metadataStore Storer
	objectStore   StorerGetter
}

type Storer interface {
	Store(context.Context, string, []byte) error
}

type Getter interface {
	Get(context.Context, string) (io.ReadCloser, error)
}

type StorerGetter interface {
	Storer
	Getter
}

func New(metadataStore Storer, objectStore StorerGetter) *Server {
	return &Server{metadataStore, objectStore}
}

func (s *Server) Store(ctx context.Context, r io.Reader) (archivistapi.StoreResponse, error) {
	payload, err := io.ReadAll(r)
	if err != nil {
		return archivistapi.StoreResponse{}, err
	}

	gid, err := gitoid.New(bytes.NewReader(payload), gitoid.WithContentLength(int64(len(payload))), gitoid.WithSha256())
	if err != nil {
		log.FromContext(ctx).Errorf("failed to generate gitoid: %v", err)
		return archivistapi.StoreResponse{}, err
	}

	if err := s.metadataStore.Store(ctx, gid.String(), payload); err != nil {
		log.FromContext(ctx).Errorf("received error from metadata store: %+v", err)
		return archivistapi.StoreResponse{}, err
	}

	if s.objectStore != nil {
		if err := s.objectStore.Store(ctx, gid.String(), payload); err != nil {
			log.FromContext(ctx).Errorf("received error from object store: %+v", err)
			return archivistapi.StoreResponse{}, err
		}
	}

	return archivistapi.StoreResponse{Gitoid: gid.String()}, nil
}

func (s *Server) StoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	resp, err := s.Store(r.Context(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		log.FromContext(r.Context()).Errorf("failed to copy storeresponse to response: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *Server) Get(ctx context.Context, gitoid string) (io.ReadCloser, error) {
	if len(strings.TrimSpace(gitoid)) == 0 {
		return nil, errors.New("gitoid parameter is required")
	}

	if s.objectStore == nil {
		return nil, errors.New("object store unavailable")
	}

	objReader, err := s.objectStore.Get(ctx, gitoid)
	if err != nil {
		log.FromContext(ctx).Errorf("failed to get object: %+v", err)
	}

	return objReader, err
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	attestationReader, err := s.Get(r.Context(), vars["gitoid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer attestationReader.Close()
	if _, err := io.Copy(w, attestationReader); err != nil {
		log.FromContext(r.Context()).Errorf("failed to copy attestation to response: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
