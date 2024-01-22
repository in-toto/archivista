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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/in-toto/archivista/pkg/api"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT APIGraphQL
type UTAPIGraphQLSuite struct {
	suite.Suite
}

func TestAPIGraphQLSuite(t *testing.T) {
	suite.Run(t, new(UTAPIGraphQLSuite))
}

func (ut *UTAPIGraphQLSuite) Test_Store() {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"data": {"data": "test"}}`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	type testSubjectVar struct {
		Gitoid string `json:"gitoid"`
	}

	type testSubjectResult struct {
		Data string `json:"data"`
	}
	result, err := api.GraphQlQuery[testSubjectResult](ctx, testServer.URL, `query`, testSubjectVar{Gitoid: "test_Gitoid"})
	ut.NoError(err)
	ut.Equal(testSubjectResult{Data: "test"}, result)
}

func (ut *UTAPIGraphQLSuite) Test_Store_NoServer() {
	ctx := context.TODO()

	type testSubjectVar struct {
		Gitoid string `json:"gitoid"`
	}

	type testSubjectResult struct {
		Data string `json:"data"`
	}
	result, err := api.GraphQlQuery[testSubjectResult](
		ctx,
		"http://invalid-archivista",
		`query`,
		testSubjectVar{Gitoid: "test_Gitoid"},
	)
	ut.Error(err)
	ut.Equal(testSubjectResult{Data: ""}, result)
}

func (ut *UTAPIGraphQLSuite) Test_Store_BadStatusCode_NoMsg() {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	type testSubjectVar struct {
		Gitoid string `json:"gitoid"`
	}

	type testSubjectResult struct {
		Data string `json:"data"`
	}
	result, err := api.GraphQlQuery[testSubjectResult](ctx, testServer.URL, `query`, testSubjectVar{Gitoid: "test_Gitoid"})
	ut.Error(err)
	ut.Equal(testSubjectResult{Data: ""}, result)
}

func (ut *UTAPIGraphQLSuite) Test_Store_InvalidData() {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(``))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	type testSubjectVar struct {
		Gitoid string `json:"gitoid"`
	}

	type testSubjectResult struct {
		Data string `json:"data"`
	}
	result, err := api.GraphQlQuery[testSubjectResult](ctx, testServer.URL, `query`, testSubjectVar{Gitoid: "test_Gitoid"})
	ut.Error(err)
	ut.Equal(testSubjectResult{Data: ""}, result)
}

func (ut *UTAPIGraphQLSuite) Test_Store_QLReponseWithErrors() {
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"data": {"data": "test"}, "errors": [{"message": "test_error"}]}}`))
				if err != nil {
					ut.FailNow(err.Error())
				}
			},
		),
	)
	defer testServer.Close()
	ctx := context.TODO()

	type testSubjectVar struct {
		Gitoid string `json:"gitoid"`
	}

	type testSubjectResult struct {
		Data   string `json:"data"`
		Errors string `json:"errors"`
	}

	result, err := api.GraphQlQuery[testSubjectResult](ctx, testServer.URL, `query`, testSubjectVar{Gitoid: "test_Gitoid"})
	ut.Error(err)
	ut.EqualError(err, "graph ql query failed: [{test_error}]")
	ut.Equal(testSubjectResult{Data: ""}, result)
}
