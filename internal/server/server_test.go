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
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/in-toto/archivista/internal/config"
	"github.com/in-toto/archivista/pkg/api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StorerMock struct {
	mock.Mock
	Storer
}

type StorerGetterMock struct {
	mock.Mock
	StorerGetter
}

type UTServerSuite struct {
	suite.Suite
	mockedStorer       *StorerMock
	mockedStorerGetter *StorerGetterMock
	testServer         Server
}

func TestUTServerSuite(t *testing.T) {
	suite.Run(t, new(UTServerSuite))
}

func (ut *UTServerSuite) SetupTest() {
	ut.mockedStorer = new(StorerMock)
	ut.mockedStorerGetter = new(StorerGetterMock)
	ut.testServer = Server{ut.mockedStorer, ut.mockedStorerGetter, nil}
}

func (ut *UTServerSuite) Test_New() {
	cfg := new(config.Config)
	cfg.EnableGraphql = true
	cfg.GraphqlWebClientEnable = true
	ut.testServer = New(ut.mockedStorer, ut.mockedStorerGetter, cfg, nil)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
	ut.testServer = New(ut.mockedStorer, ut.mockedStorerGetter, cfg, nil)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
	ut.testServer = New(ut.mockedStorer, ut.mockedStorerGetter, cfg, nil)
	ut.NotNil(ut.testServer)
	router := ut.testServer.Router()
	ut.NotNil(router)

	allPaths := []string{}
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

// Mock StorerGetter.Store()
func (m *StorerGetterMock) Store(context.Context, string, []byte) error {
	args := m.Called()
	return args.Error(0)
}

// Mock StorerGetter.Get()
func (m *StorerGetterMock) Get(context.Context, string) (io.ReadCloser, error) {
	args := m.Called()
	stringReader := strings.NewReader("testData")
	stringReadCloser := io.NopCloser(stringReader)
	return stringReadCloser, args.Error(0)
}

// Mock StorerMock.Store()
func (m *StorerMock) Store(context.Context, string, []byte) error {
	args := m.Called()
	return args.Error(0)
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

func (ut *UTServerSuite) Test_Download_EmptyGitoid() {
	ctx := context.TODO()
	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	_, err := ut.testServer.Download(ctx, "")
	ut.ErrorContains(err, "gitoid parameter is required")
}

func (ut *UTServerSuite) Test_Download_EmptyGitoidTrimmed() {
	ctx := context.TODO()
	ut.mockedStorerGetter.On("Get").Return(nil) // mock Get() to return nil

	_, err := ut.testServer.Download(ctx, "           ")
	ut.ErrorContains(err, "gitoid parameter is required")
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

func (ut *UTServerSuite) Test_DownloadHandler_ObjectStorageFailed() {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/download", nil)
	request = mux.SetURLVars(request, map[string]string{"gitoid": "fakeGitoid"})

	ut.mockedStorerGetter.On("Get").Return(errors.New("BAD S3")) // mock Get() to return nil

	ut.testServer.DownloadHandler(w, request)
	ut.Equal(http.StatusInternalServerError, w.Code)
	ut.Contains(w.Body.String(), "BAD S3")
}
