// Copyright 2022-2024 The Archivista Contributors
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

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edwarnicke/gitoid"
	"github.com/gorilla/mux"
	"github.com/in-toto/archivista"
	_ "github.com/in-toto/archivista/docs"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/internal/config"
	"github.com/in-toto/archivista/pkg/api"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	metadataStore Storer
	objectStore   StorerGetter
	router        *mux.Router
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

func New(metadataStore Storer, objectStore StorerGetter, cfg *config.Config, sqlClient *ent.Client) Server {
	r := mux.NewRouter()
	s := &Server{metadataStore, objectStore, nil}

	// TODO: remove from future version (v0.6.0) endpoint with version
	r.HandleFunc("/download/{gitoid}", s.DownloadHandler)
	r.HandleFunc("/upload", s.UploadHandler)
	if cfg.EnableGraphql {
		r.Handle("/query", s.Query(sqlClient))
		r.Handle("/v1/query", s.Query(sqlClient))
	}

	r.HandleFunc("/v1/download/{gitoid}", s.DownloadHandler)
	r.HandleFunc("/v1/upload", s.UploadHandler)
	if cfg.GraphqlWebClientEnable {
		r.Handle("/",
			playground.Handler("Archivista", "/v1/query"),
		)
	}
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	s.router = r

	return *s

}

// @title Archivista API
// @description Archivista API
// @version v1
// @contact.name Archivista Contributors
// @contact.url https://github.com/in-toto/archivista/issues/new
// @license Apache 2
// @license.url https://opensource.org/licenses/Apache-2
// InitRoutes initializes the HTTP API routes for the server
func (s *Server) Router() *mux.Router {
	return s.router
}

// @Summary Upload
// @Description stores an attestation
// @Produce  json
// @Success 200 {object} api.StoreResponse
// @Tags attestation
// @Router /v1/upload [post]
func (s *Server) Upload(ctx context.Context, r io.Reader) (api.UploadResponse, error) {
	payload, err := io.ReadAll(r)
	if err != nil {
		return api.UploadResponse{}, err
	}

	gid, err := gitoid.New(bytes.NewReader(payload), gitoid.WithContentLength(int64(len(payload))), gitoid.WithSha256())
	if err != nil {
		logrus.Errorf("failed to generate gitoid: %v", err)
		return api.UploadResponse{}, err
	}

	if s.objectStore != nil {
		if err := s.objectStore.Store(ctx, gid.String(), payload); err != nil {
			logrus.Errorf("received error from object store: %+v", err)
			return api.UploadResponse{}, err
		}
	}

	if err := s.metadataStore.Store(ctx, gid.String(), payload); err != nil {
		logrus.Errorf("received error from metadata store: %+v", err)
		return api.UploadResponse{}, err
	}

	return api.UploadResponse{Gitoid: gid.String()}, nil
}

// @Summary Upload
// @Description stores an attestation
// @Produce  json
// @Success 200 {object} api.StoreResponse
// @Router /upload [post]
// @Deprecated
func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	resp, err := s.Upload(r.Context(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		logrus.Errorf("failed to copy storeresponse to response: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// @Summary Download
// @Description download an attestation
// @Produce  json
// @Success 200 {object} dsse.Envelope
// @Tags attestation
// @Router /v1/download/{gitoid} [post]
func (s *Server) Download(ctx context.Context, gitoid string) (io.ReadCloser, error) {
	if len(strings.TrimSpace(gitoid)) == 0 {
		return nil, errors.New("gitoid parameter is required")
	}

	if s.objectStore == nil {
		return nil, errors.New("object store unavailable")
	}

	objReader, err := s.objectStore.Get(ctx, gitoid)
	if err != nil {
		logrus.Errorf("failed to get object: %+v", err)
	}

	return objReader, err
}

// @Summary Download
// @Description download an attestation
// @Produce  json
// @Success 200 {object} dsse.Envelope
// @Deprecated
// @Router /download/{gitoid} [post]
func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	attestationReader, err := s.Download(r.Context(), vars["gitoid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer attestationReader.Close()
	if _, err := io.Copy(w, attestationReader); err != nil {
		logrus.Errorf("failed to copy attestation to response: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// @Summary Query GraphQL
// @Description GraphQL query
// @Produce  json
// @Success 200 {object} archivista.Resolver
// @Tags graphql
// @Router /v1/query [post]
func (s *Server) Query(sqlclient *ent.Client) *handler.Server {
	srv := handler.NewDefaultServer(archivista.NewSchema(sqlclient))
	srv.Use(entgql.Transactioner{TxOpener: sqlclient})
	return srv
}
