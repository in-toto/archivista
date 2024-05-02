// Copyright 2023 The Archivista Contributors
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

package api_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/in-toto/go-witness/dsse"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT APIDownload
type UTAPIDownloadSuite struct {
	suite.Suite
}

func TestAPIDownloadSuite(t *testing.T) {
	suite.Run(t, new(UTAPIDownloadSuite))
}

func (ut *UTAPIDownloadSuite) Test_Download() {

	testEnvelope, err := os.ReadFile("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}
	expectedEnvelop := dsse.Envelope{}
	err = json.Unmarshal(testEnvelope, &expectedEnvelop)
	if err != nil {
		ut.FailNow(err.Error())
	}

	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(testEnvelope)
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	// test api.Download happy flow
	resp, err := api.Download(ctx, testServer.URL, "gitoid_test")
	if err != nil {
		ut.FailNow(err.Error())
	}

	ut.Equal(expectedEnvelop, resp)
}

func (ut *UTAPIDownloadSuite) Test_Download_DecodeFailure() {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`invalid`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	// test api.Download happy flow
	resp, err := api.Download(ctx, testServer.URL, "gitoid_test")
	ut.Error(err)
	ut.Equal(dsse.Envelope{}, resp)
}

func (ut *UTAPIDownloadSuite) Test_DownloadWithReader() {

	// mock server
	expected := `{"body":"body"}`
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(expected))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()

	// context
	ctx := context.TODO()

	// temporary file
	tempDir := os.TempDir()
	dst, err := os.Create(path.Join(tempDir, "test"))
	if err != nil {
		ut.FailNow(err.Error())
	}
	err = api.DownloadWithWriter(ctx, testServer.URL, "gitoid", dst)
	if err != nil {
		ut.FailNow(err.Error())
	}

	// validate result
	result, err := os.ReadFile(path.Join(tempDir, "test"))
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(expected, string(result))
}

func (ut *UTAPIDownloadSuite) Test_DownloadWithWriter_NoServer() {

	// context
	ctx := context.TODO()

	// dst as stdout
	var dst io.Writer = os.Stdout

	err := api.DownloadWithWriter(ctx, "http://invalid-archivista", "gitoid_test", dst)
	ut.Error(err)
}

func (ut *UTAPIDownloadSuite) Test_DownloadWithWriter_BadStatusCode() {

	// mock server
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(`Internal Server Error`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()

	// dst as stdout
	var dst io.Writer = os.Stdout

	// context
	ctx := context.TODO()

	err := api.DownloadWithWriter(ctx, testServer.URL, "gitoid_test", dst)
	ut.ErrorContains(err, "Internal Server Error")
}
