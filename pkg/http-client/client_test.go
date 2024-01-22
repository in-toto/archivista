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

package httpclient_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/in-toto/archivista/pkg/api"
	httpclient "github.com/in-toto/archivista/pkg/http-client"
	"github.com/in-toto/go-witness/dsse"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT HTTPClientDownloadSuite
type UTHTTPClientDownloadSuite struct {
	suite.Suite
}

func TestHTTPClientAPIDownloadSuite(t *testing.T) {
	suite.Run(t, new(UTHTTPClientDownloadSuite))
}

func (ut *UTHTTPClientDownloadSuite) Test_DownloadDSSE() {
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
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	resp, err := client.DownloadDSSE(ctx, "gitoid_test")
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(expectedEnvelop, resp)
}

func (ut *UTHTTPClientDownloadSuite) Test_DownloadReadCloser() {
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
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	readCloser, err := client.DownloadReadCloser(ctx, "gitoid_test")
	if err != nil {
		ut.FailNow(err.Error())
	}
	env := dsse.Envelope{}
	if err := json.NewDecoder(readCloser).Decode(&env); err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(expectedEnvelop, env)
}

func (ut *UTHTTPClientDownloadSuite) Test_DownloadWithWriter() {
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
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	buf := bytes.NewBuffer(nil)
	if err := client.DownloadWithWriter(ctx, "gitoid_test", buf); err != nil {
		ut.FailNow(err.Error())
	}
	env := dsse.Envelope{}
	if err := json.NewDecoder(buf).Decode(&env); err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(expectedEnvelop, env)
}

// Test Suite: UT HTTPClientStore
type UTHTTPClientStoreSuite struct {
	suite.Suite
}

func TestAPIStoreSuite(t *testing.T) {
	suite.Run(t, new(UTHTTPClientStoreSuite))
}

func (ut *UTHTTPClientStoreSuite) Test_Store() {
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
	attFile, err := os.ReadFile("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}
	attEnvelop := dsse.Envelope{}
	err = json.Unmarshal(attFile, &attEnvelop)
	if err != nil {
		ut.FailNow(err.Error())
	}
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	resp, err := client.Store(ctx, attEnvelop)
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(resp, api.StoreResponse{Gitoid: "test"})
}

func (ut *UTHTTPClientStoreSuite) Test_StoreWithReader() {
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
	attIo, err := os.Open("../../test/package.attestation.json")
	if err != nil {
		ut.FailNow(err.Error())
	}
	ctx := context.TODO()
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	resp, err := client.StoreWithReader(ctx, attIo)
	if err != nil {
		ut.FailNow(err.Error())
	}
	ut.Equal(resp, api.StoreResponse{Gitoid: "test"})
}

// Test Suite: UT HTTPClientStore
type UTHTTPClientGraphQLSuite struct {
	suite.Suite
}

func TestAPIGraphQLSuite(t *testing.T) {
	suite.Run(t, new(UTHTTPClientGraphQLSuite))
}

func (ut *UTHTTPClientGraphQLSuite) Test_GraphQLRetrieveSubjectResults() {
	expected := api.GraphQLResponseGeneric[api.RetrieveSubjectResults]{
		Data: api.RetrieveSubjectResults{
			Subjects: api.Subjects{
				Edges: []api.SubjectEdge{
					{
						Node: api.SubjectNode{
							Name: "test_Gitoid",
							SubjectDigests: []api.SubjectDigest{
								{
									Algorithm: "test_Gitoid",
									Value:     "test_Gitoid",
								},
							},
						},
					},
				},
			},
		},
		Errors: []api.GraphQLError{},
	}
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(expected); err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	actual, err := client.GraphQLRetrieveSubjectResults(ctx, "test_Gitoid")
	ut.NoError(err)
	ut.Equal(expected.Data, actual)
}

func (ut *UTHTTPClientGraphQLSuite) Test_GraphQLSearchResults() {
	expected := api.GraphQLResponseGeneric[api.SearchResults]{
		Data: api.SearchResults{
			Dsses: api.DSSES{
				Edges: []api.SearchEdge{
					{
						Node: api.SearchNode{
							GitoidSha256: "test_Gitoid",
							Statement: api.Statement{
								AttestationCollection: api.AttestationCollection{
									Name: "test_Gitoid",
									Attestations: []api.Attestation{
										{
											Type: "test",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Errors: []api.GraphQLError{},
	}
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(expected); err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()
	client, err := httpclient.CreateArchivistaClient(http.DefaultClient, testServer.URL)
	if err != nil {
		ut.FailNow(err.Error())
	}
	actual, err := client.GraphQLRetrieveSearchResults(ctx, "test_Gitoid", "test_Gitoid")
	ut.NoError(err)
	ut.Equal(expected.Data, actual)
}
