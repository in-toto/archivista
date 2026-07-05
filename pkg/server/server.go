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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/edwarnicke/gitoid"
	"github.com/gorilla/mux"
	"github.com/in-toto/archivista"
	_ "github.com/in-toto/archivista/docs"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/api"
	"github.com/in-toto/archivista/pkg/artifactstore"
	"github.com/in-toto/archivista/pkg/config"
	"github.com/in-toto/archivista/pkg/publisherstore"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	metadataStore  Storer
	objectStore    StorerGetter
	artifactStore  artifactstore.Store
	router         *mux.Router
	sqlClient      *ent.Client
	publisherStore []publisherstore.Publisher
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

type Option func(*Server)

func WithMetadataStore(metadataStore Storer) Option {
	return func(s *Server) {
		s.metadataStore = metadataStore
	}
}

func WithObjectStore(objectStore StorerGetter) Option {
	return func(s *Server) {
		s.objectStore = objectStore
	}
}

func WithEntSqlClient(sqlClient *ent.Client) Option {
	return func(s *Server) {
		s.sqlClient = sqlClient
	}
}

func WithArtifactStore(wds artifactstore.Store) Option {
	return func(s *Server) {
		s.artifactStore = wds
	}
}

func WithPublishers(pub []publisherstore.Publisher) Option {
	return func(s *Server) {
		s.publisherStore = pub
	}
}

func New(cfg *config.Config, opts ...Option) (Server, error) {
	r := mux.NewRouter()
	s := Server{
		router: r,
	}

	for _, opt := range opts {
		opt(&s)
	}

	// TODO: remove from future version (v0.6.0) endpoint with version
	r.HandleFunc("/download/{gitoid}", s.DownloadHandler)
	r.HandleFunc("/upload", s.UploadHandler)
	if cfg.EnableSQLStore && cfg.EnableGraphql {
		r.Handle("/query", s.Query(s.sqlClient))
		r.Handle("/v1/query", s.Query(s.sqlClient))
	}

	r.HandleFunc("/v1/download/{gitoid}", s.DownloadHandler)
	r.HandleFunc("/v1/upload", s.UploadHandler)
	if cfg.EnableSQLStore && cfg.EnableGraphql && cfg.GraphqlWebClientEnable {
		r.Handle("/",
			playground.Handler("Archivista", "/v1/query"),
		)
	}

	if cfg.EnableArtifactStore {
		r.HandleFunc("/v1/artifacts", s.AllArtifactsHandler)
		r.HandleFunc("/v1/artifacts/{name}", s.ArtifactAllVersionsHandler)
		r.HandleFunc("/v1/artifacts/{name}/{version}", s.ArtifactVersionHandler)
		r.HandleFunc("/v1/download/artifact/{name}/{version}/{distribution}", s.DownloadArtifactHandler)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return s, nil
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

	if s.metadataStore != nil {
		if err := s.metadataStore.Store(ctx, gid.String(), payload); err != nil {
			logrus.Errorf("received error from metadata store: %+v", err)
			return api.UploadResponse{}, err
		}
	}

	if s.publisherStore != nil {
		for _, publisher := range s.publisherStore {
			// TODO: Make publish asynchrouns and use goroutine
			if err := publisher.Publish(ctx, gid.String(), payload); err != nil {
				logrus.Errorf("received error from publisher: %+v", err)
			}
		}
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
// @Param gitoid path string true "gitoid"
// @Success 200 {object} dsse.Envelope
// @Failure 500 {object} string
// @Failure 404 {object} nil
// @Failure 400 {object} string
// @Tags attestation
// @Router /v1/download/{gitoid} [get]
func (s *Server) Download(ctx context.Context, gitoid string) (io.ReadCloser, error) {
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
// @Param gitoid path string true "gitoid"
// @Success 200 {object} dsse.Envelope
// @Failure 500 {object} string
// @Failure 404 {object} nil
// @Failure 400 {object} string
// @Deprecated
// @Router /download/{gitoid} [get]
func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "gitoid parameter is required", http.StatusBadRequest)
		return
	}
	if len(strings.TrimSpace(vars["gitoid"])) == 0 {
		http.Error(w, "gitoid parameter is required", http.StatusBadRequest)
		return
	}

	attestationReader, err := s.Download(r.Context(), vars["gitoid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer attestationReader.Close()
	if _, err := io.Copy(w, attestationReader); err != nil {
		logrus.Errorf("failed to copy attestation to response: %+v", err)
		w.WriteHeader(http.StatusNotFound)
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
	srv := handler.New(archivista.NewSchema(sqlclient))
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(entgql.Transactioner{TxOpener: sqlclient})
	return srv
}

// @Summary List all Artifacts
// @Description retrieves details about all available Artifacts
// @Produce json
// @Success 200 {object} map[string]artifactstore.Artifact
// @Failure 500 {object} string
// @Failure 400 {object} string
// @Tags Artifacts
// @Router /v1/artifacts [get]
func (s *Server) AllArtifactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	allArtifacts := s.artifactStore.Artifacts()
	allArtifactsJson, err := json.Marshal(allArtifacts)
	if err != nil {
		http.Error(w, fmt.Errorf("could not marshal artifact versions: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(allArtifactsJson)); err != nil {
		http.Error(w, fmt.Errorf("could not send json response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// @Summary List Artifact Versions
// @Description retrieves details about all available versions of a specified artifact
// @Produce json
// @Param name path string true "artifact name"
// @Success 200 {object} map[string]artifactstore.Version
// @Failure 500 {object} string
// @Failure 400 {object} string
// @Tags Artifacts
// @Router /v1/artifacts/{name} [get]
func (s *Server) ArtifactAllVersionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	artifactName := vars["name"]
	if len(artifactName) == 0 {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	artifactVersions, ok := s.artifactStore.Versions(artifactName)
	if !ok {
		http.Error(w, "artifact not found", http.StatusNotFound)
		return
	}

	artifactVersionsJson, err := json.Marshal(artifactVersions)
	if err != nil {
		http.Error(w, fmt.Errorf("could not marshal artifact versions: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(artifactVersionsJson)); err != nil {
		http.Error(w, fmt.Errorf("could not send json response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// @Summary Artifact Version Details
// @Description retrieves details about a specified version of an artifact
// @Produce json
// @Param name path string true "artifact name"
// @Param version path string true "version of artifact"
// @Success 200 {object} artifactstore.Version
// @Failure 500 {objecpec} string
// @Failure 404 {object} nil
// @Failure 400 {object} string
// @Tags Artifacts
// @Router /v1/artifacts/{name}/{version} [get]
func (s *Server) ArtifactVersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "name and version parameters are required", http.StatusBadRequest)
		return
	}

	artifactString := vars["name"]
	if len(artifactString) == 0 {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	versionString := vars["version"]
	if len(versionString) == 0 {
		http.Error(w, "version parameter is required", http.StatusBadRequest)
		return
	}

	version, ok := s.artifactStore.Version(artifactString, versionString)
	if !ok {
		http.Error(w, "version not found", http.StatusNotFound)
		return
	}

	versionJson, err := json.Marshal(version)
	if err != nil {
		http.Error(w, fmt.Errorf("could not marshal artifact distros: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, bytes.NewReader(versionJson)); err != nil {
		http.Error(w, fmt.Errorf("could not send json response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// @Summary Download Artifact
// @Description downloads a specified distribution of an artifact
// @Produce octet-stream
// @Param name path string true "name of artifact"
// @Param version path string true "version of artifact to download"
// @Param distribution path string true "distribution of artifact to download"
// @Success 200 {file} octet-stream
// @Failure 500 {object} string
// @Failure 404 {object} nil
// @Failure 400 {object} string
// @Tags Artifacts
// @Router /v1/download/artifact/{name}/{version}/{distribution} [get]
func (s *Server) DownloadArtifactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s is an unsupported method", r.Method), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	if vars == nil {
		http.Error(w, "version and distribution parameter is required", http.StatusBadRequest)
		return
	}

	artifactString := vars["name"]
	if len(artifactString) == 0 {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	versionString := vars["version"]
	if len(versionString) == 0 {
		http.Error(w, "version parameter is required", http.StatusBadRequest)
		return
	}

	distroString := vars["distribution"]
	if len(distroString) == 0 {
		http.Error(w, "distribution parameter is required", http.StatusBadRequest)
		return
	}

	distro, ok := s.artifactStore.Distribution(artifactString, versionString, distroString)
	if !ok {
		http.Error(w, "distribution of artifact not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(distro.FileLocation)
	if err != nil {
		http.Error(w, "could not read distribution file", http.StatusBadRequest)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			logrus.Errorf("failed to close artifact distribution file %s: %+v", distro.FileLocation, err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "could not stat distribution file", http.StatusBadRequest)
		return
	}

	fileName := filepath.Base(distro.FileLocation)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, fmt.Errorf("could not send artifact distribution: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}
