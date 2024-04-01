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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/in-toto/go-witness/dsse"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT APIStore
type UTAPIStoreSuite struct {
	suite.Suite
}

func TestAPIStoreSuite(t *testing.T) {
	suite.Run(t, new(UTAPIStoreSuite))
}

func (ut *UTAPIStoreSuite) Test_Store() {

	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"gitoid":"test"}`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	// load test valid test file and parse the dsse envelop
	attFile, err := os.ReadFile("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}
	attEnvelop := &dsse.Envelope{}
	err = json.Unmarshal(attFile, &attEnvelop)
	if err != nil {
		ut.FailNow(err.Error())
	}

	// test api.Store happy flow
	resp, err := api.Store(ctx, testServer.URL, *attEnvelop)
	if err != nil {
		ut.FailNow(err.Error())
	}

	ut.Equal(resp, api.StoreResponse{Gitoid: "test"})
}

func (ut *UTAPIStoreSuite) Test_StoreWithReader() {

	// mock server
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"gitoid":"test"}`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()

	// io.Reader file
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// context
	ctx := context.TODO()

	resp, err := api.StoreWithReader(ctx, testServer.URL, attIo)
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(resp, api.StoreResponse{Gitoid: "test"})
}

func (ut *UTAPIStoreSuite) Test_StoreWithReader_NoServer() {

	// io.Reader file
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// context
	ctx := context.TODO()

	resp, err := api.StoreWithReader(ctx, "http://invalid-archivista", attIo)
	ut.Error(err)
	ut.Equal(resp, api.StoreResponse{})
}

func (ut *UTAPIStoreSuite) Test_StoreWithReader_InvalidResponseBody() {

	// mock server
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"invalid":"invalid"}`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()

	// io.Reader file
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// context
	ctx := context.TODO()

	resp, err := api.StoreWithReader(ctx, testServer.URL, attIo)
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(resp, api.StoreResponse{})
}

func (ut *UTAPIStoreSuite) Test_StoreWithReader_BadStatusCode() {

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

	// io.Reader file
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// context
	ctx := context.TODO()

	resp, err := api.StoreWithReader(ctx, testServer.URL, attIo)
	ut.ErrorContains(err, "Internal Server Error")
	ut.Equal(resp, api.StoreResponse{})
}

func (ut *UTAPIStoreSuite) Test_StoreWithReader_BadJSONBody() {

	// mock server
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	defer testServer.Close()

	// io.Reader file
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}

	// context
	ctx := context.TODO()

	resp, err := api.StoreWithReader(ctx, testServer.URL, attIo)
	ut.ErrorContains(err, "unexpected end of JSON input")
	ut.Equal(resp, api.StoreResponse{})
}
