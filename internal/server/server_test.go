// Copyright 2024 The Archivista Contributors
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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/in-toto/archivista/internal/artifactstore"
	"github.com/in-toto/archivista/internal/config"
	"github.com/in-toto/archivista/pkg/api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock StorerMock
type StorerMock struct {
	mock.Mock
	Storer
}

func (m *StorerMock) Store(context.Context, string, []byte) error {
	args := m.Called()
	return args.Error(0)
}

type StorerGetterMock struct {
	mock.Mock
	StorerGetter
}

// Mock StorerGetter
func (m *StorerGetterMock) Store(context.Context, string, []byte) error {
	args := m.Called()
	return args.Error(0)
}

func (m *StorerGetterMock) Get(context.Context, string) (io.ReadCloser, error) {
	args := m.Called()
	stringReader := strings.NewReader("testData")
	stringReadCloser := io.NopCloser(stringReader)
	return stringReadCloser, args.Error(0)
}

// Mock ResponseRecorderMock
type ResponseRecorderMock struct {
	mock.Mock
	Code      int
	HeaderMap http.Header
	Body      *bytes.Buffer
	Flushed   bool
}

func (m *ResponseRecorderMock) Write([]byte) (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (rw *ResponseRecorderMock) Header() http.Header {
	return http.Header{}
}

func (rw *ResponseRecorderMock) WriteHeader(code int) {
	rw.Code = code
}

func (rw *ResponseRecorderMock) WriteString(str string) (int, error) {
	return 0, nil
}

// Define test Suite
type UTServerSuite struct {
	suite.Suite
	mockedStorer          *StorerMock
	mockedStorerGetter    *StorerGetterMock
	mockedResposeRecorder *ResponseRecorderMock
	testServer            Server
}

func TestUTServerSuite(t *testing.T) {
	suite.Run(t, new(UTServerSuite))
}

func (ut *UTServerSuite) SetupTest() {
	ut.mockedStorer = new(StorerMock)
	ut.mockedStorerGetter = new(StorerGetterMock)
	ut.mockedResposeRecorder = new(ResponseRecorderMock)
	ut.testServer = Server{
		metadataStore: ut.mockedStorer,
		objectStore:   ut.mockedStorerGetter,
		artifactStore: ut.testArtifactStore(),
	}
}

func (ut *UTServerSuite) Test_New() {
	cfg := new(config.Config)
	cfg.EnableGraphql = true
	cfg.GraphqlWebClientEnable = true
	var err error
	ut.testServer, err = New(cfg, WithMetadataStore(ut.mockedStorer), WithObjectStore(ut.mockedStorerGetter))
	ut.NoError(err)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			ut.FailNow(err.Error())
		}
		allPaths = append(allPaths, pathTemplate)
		return nil
	})

	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Contains(allPaths, "/download/{gitoid}")
	ut.Contains(allPaths, "/upload")
	ut.Contains(allPaths, "/query")
	ut.Contains(allPaths, "/v1/download/{gitoid}")
	ut.Contains(allPaths, "/v1/upload")
	ut.Contains(allPaths, "/v1/query")
	ut.Contains(allPaths, "/")
	ut.Contains(allPaths, "/swagger/")
}

func (ut *UTServerSuite) Test_New_EnableGraphQL_False() {
	cfg := new(config.Config)
	cfg.EnableGraphql = false
	cfg.GraphqlWebClientEnable = true
	var err error
	ut.testServer, err = New(cfg, WithMetadataStore(ut.mockedStorer), WithObjectStore(ut.mockedStorerGetter))
	ut.NoError(err)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			ut.FailNow(err.Error())
		}
		allPaths = append(allPaths, pathTemplate)
		return nil
	})

	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Contains(allPaths, "/download/{gitoid}")
	ut.Contains(allPaths, "/upload")
	ut.NotContains(allPaths, "/query")
	ut.Contains(allPaths, "/v1/download/{gitoid}")
	ut.Contains(allPaths, "/v1/upload")
	ut.NotContains(allPaths, "/v1/query")
	ut.Contains(allPaths, "/")
	ut.Contains(allPaths, "/swagger/")
}

func (ut *UTServerSuite) Test_New_GraphqlWebClientEnable_False() {
	cfg := new(config.Config)
	cfg.EnableGraphql = true
	cfg.GraphqlWebClientEnable = false
	var err error
	ut.testServer, err = New(cfg, WithMetadataStore(ut.mockedStorer), WithObjectStore(ut.mockedStorerGetter))
	ut.NoError(err)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			ut.FailNow(err.Error())
		}
		allPaths = append(allPaths, pathTemplate)
		return nil
	})

	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Contains(allPaths, "/download/{gitoid}")
	ut.Contains(allPaths, "/upload")
	ut.Contains(allPaths, "/query")
	ut.Contains(allPaths, "/v1/download/{gitoid}")
	ut.Contains(allPaths, "/v1/upload")
	ut.Contains(allPaths, "/v1/query")
	ut.NotContains(allPaths, "/")
	ut.Contains(allPaths, "/swagger/")
}

func (ut *UTServerSuite) Test_Upload() {
	ctx := context.TODO()
	r := strings.NewReader("fakeTestData")

	ut.mockedStorerGetter.On("Store").Return(nil) // mock Get() to return nil
	ut.mockedStorer.On("Store").Return(nil)       // mock Store() to return nil

	resp, err := ut.testServer.Upload(ctx, r)
	ut.NoError(err)
	ut.NotEqual("", resp.Gitoid)
}

func (ut *UTServerSuite) Test_Upload_NoObjectStorage() {
	ctx := context.TODO()
	r := strings.NewReader("fakeTestData")

	ut.testServer.objectStore = nil
	ut.mockedStorer.On("Store").Return(nil) // mock Store() to return nil

	resp, err := ut.testServer.Upload(ctx, r)
	ut.NoError(err)
	ut.NotEqual("", resp.Gitoid)
}

func (ut *UTServerSuite) Test_Upload_FailedObjectStorage() {
	ctx := context.TODO()
	r := strings.NewReader("fakeTestData")

	ut.mockedStorerGetter.On("Store").Return(errors.New("Bad S3")) // mock Get() to return err
	ut.mockedStorer.On("Store").Return(nil)                        // mock Store() to return nil

	resp, err := ut.testServer.Upload(ctx, r)
	ut.ErrorContains(err, "Bad S3")
	ut.Equal(api.UploadResponse{}, resp)
}

func (ut *UTServerSuite) Test_Upload_FailedMetadatStprage() {
	ctx := context.TODO()
	r := strings.NewReader("fakeTestData")

	ut.mockedStorerGetter.On("Store").Return(nil)             // mock Get() to return nil
	ut.mockedStorer.On("Store").Return(errors.New("Bad SQL")) // mock Store() to return err

	resp, err := ut.testServer.Upload(ctx, r)
	ut.ErrorContains(err, "Bad SQL")
	ut.Equal(api.UploadResponse{}, resp)
}

func (ut *UTServerSuite) Test_UploadHandler() {

	w := httptest.NewRecorder()
	requestBody := []byte("fakePayload")
	request := httptest.NewRequest(http.MethodPost, "/v1/upload", bytes.NewBuffer(requestBody))

	ut.mockedStorerGetter.On("Store").Return(nil) // mock Get() to return nil
	ut.mockedStorer.On("Store").Return(nil)       // mock Store() to return nil

	ut.testServer.UploadHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
}

func (ut *UTServerSuite) Test_UploadHandler_WrongMethod() {

	w := httptest.NewRecorder()
	requestBody := []byte("fakePayload")
	request := httptest.NewRequest(http.MethodGet, "/upload", bytes.NewBuffer(requestBody))

	ut.mockedStorerGetter.On("Store").Return(nil) // mock Get() to return nil
	ut.mockedStorer.On("Store").Return(nil)       // mock Store() to return nil

	ut.testServer.UploadHandler(w, request)
	ut.Equal(http.StatusBadRequest, w.Code)
	ut.Contains(w.Body.String(), "is an unsupported method")
}

func (ut *UTServerSuite) Test_UploadHandler_FailureUpload() {

	w := httptest.NewRecorder()
	requestBody := []byte("fakePayload")
	request := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(requestBody))

	ut.mockedStorerGetter.On("Store").Return(errors.New("BAD S3")) // mock Get() to return nil
	ut.mockedStorer.On("Store").Return(nil)                        // mock Store() to return nil

	ut.testServer.UploadHandler(w, request)
	ut.Equal(http.StatusInternalServerError, w.Code)
	ut.Contains(w.Body.String(), "BAD S3")
}

func (ut *UTServerSuite) Test_Download() {
	ctx := context.TODO()
	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	resp, err := ut.testServer.Download(ctx, "fakeGitoid")
	ut.NoError(err)
	data, _ := io.ReadAll(resp)
	ut.Equal("testData", string(data))
}

func (ut *UTServerSuite) Test_Download_NoObjectStorage() {
	ctx := context.TODO()
	ut.testServer.objectStore = nil

	_, err := ut.testServer.Download(ctx, "fakeGitoid")
	ut.ErrorContains(err, "object store unavailable")
}

func (ut *UTServerSuite) Test_Download_ObjectStorageError() {
	ctx := context.TODO()
	ut.mockedStorerGetter.On("Get").Return(errors.New("BAD S3")) // mock Get() to return nil

	_, err := ut.testServer.Download(ctx, "fakeGitoid")
	ut.ErrorContains(err, "BAD S3")
}

func (ut *UTServerSuite) Test_DownloadHandler() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": "fakeGitoid"})

	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
	ut.Equal("testData", w.Body.String())
}

func (ut *UTServerSuite) Test_DownloadHandler_BadMethod() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": "fakeGitoid"})

	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusBadRequest, w.Code)
	ut.Contains(w.Body.String(), "POST is an unsupported method")
}

func (ut *UTServerSuite) Test_DownloadHandler_MissingGitOID() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)

	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusBadRequest, w.Code)
	ut.Contains(w.Body.String(), "gitoid parameter is required")
}

func (ut *UTServerSuite) Test_DownloadHandler_GitOIDEmpty() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": ""})

	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusBadRequest, w.Code)
	ut.Contains(w.Body.String(), "gitoid parameter is required")
}

func (ut *UTServerSuite) Test_DownloadHandler_ObjectStorageFailed() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": "fakeGitoid"})

	ut.mockedStorerGetter.On("Get").Return(errors.New("BAD S3")) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusInternalServerError, w.Code)
	ut.Contains(w.Body.String(), "BAD S3")
}

func (ut *UTServerSuite) Test_DownloadHandler_NotFound() {
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": "fakeGitoid"})

	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil
	ut.mockedResposeRecorder.On("Write").Return(404, errors.New("Not Found"))

	ut.testServer.DownloadHandler(ut.mockedResposeRecorder, request)
	ut.Equal(http.StatusNotFound, ut.mockedResposeRecorder.Code)
	ut.Nil(ut.mockedResposeRecorder.Body)
}

func (ut *UTServerSuite) Test_AllArtifactsHandler() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/artifacts", nil)

	ut.testServer.AllArtifactsHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
	ut.Contains(w.Body.String(), "witness")
	ut.Contains(w.Body.String(), "v0.1.0")
	ut.Contains(w.Body.String(), "linux-x64")
}

func (ut *UTServerSuite) Test_ArtifactAllVersionsHandler() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/artifacts/witness", nil)
	request = mux.SetURLVars(request, map[string]string{"artifact": "witness"})

	ut.testServer.ArtifactAllVersionsHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
	ut.Contains(w.Body.String(), "v0.1.0")
	ut.Contains(w.Body.String(), "linux-x64")
}

func (ut *UTServerSuite) Test_ArtifactVersionHandler() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/artifacts/witness/v0.1.0", nil)
	request = mux.SetURLVars(request, map[string]string{"artifact": "witness", "version": "v0.1.0"})

	ut.testServer.ArtifactVersionHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
	ut.Contains(w.Body.String(), "linux-x64")
}

func (ut *UTServerSuite) Test_ArtifactVersionHandler_NotFound() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/artifacts/witness/v0.3.0", nil)
	request = mux.SetURLVars(request, map[string]string{"artifact": "witness", "version": "v0.3.0"})

	ut.testServer.ArtifactVersionHandler(w, request)
	ut.Equal(http.StatusNotFound, w.Code)
	ut.Contains(w.Body.String(), "version not found")
}

func (ut *UTServerSuite) Test_DownloadArtifactHandler() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download/artifact/witness/v0.1.0/linux-x64", nil)
	request = mux.SetURLVars(request, map[string]string{"artifact": "witness", "version": "v0.1.0", "distribution": "linux-x64"})

	ut.testServer.DownloadArtifactHandler(w, request)
	ut.Equal(http.StatusOK, w.Code)
	ut.Contains(w.Body.String(), "test")
}

func (ut *UTServerSuite) Test_DownloadArtifactHandler_NotFound() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download/artifact/witness/v0.1.0/linux-arm", nil)
	request = mux.SetURLVars(request, map[string]string{"artifact": "witness", "version": "v0.1.0", "distribution": "linux-arm"})

	ut.testServer.DownloadArtifactHandler(w, request)
	ut.Equal(http.StatusNotFound, w.Code)
	ut.Contains(w.Body.String(), "distribution of artifact not found")
}

func (ut *UTServerSuite) testArtifactStore() artifactstore.Store {
	testDir := ut.T().TempDir()
	testDistroFilePath := filepath.Join(testDir, "witness-v0.1.0-linux-x64")
	ut.NoError(os.WriteFile(testDistroFilePath, []byte("test"), 0644))

	config := artifactstore.Config{
		Artifacts: map[string]artifactstore.Artifact{
			"witness": {
				Versions: map[string]artifactstore.Version{
					"v0.1.0": {
						Description: "some description",
						Distributions: map[string]artifactstore.Distribution{
							"linux-x64": {
								SHA256Digest: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
								FileLocation: testDistroFilePath,
							},
						},
					},
				},
			},
		},
	}

	wds, err := artifactstore.New(artifactstore.WithConfig(config))
	ut.NoError(err)
	return wds
}
