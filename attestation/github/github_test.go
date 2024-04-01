// Copyright 2021 The Witness Contributors
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

package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockServer() *httptest.Server {
	type Response struct {
		Count int    `json:"count"`
		Value string `json:"value"`
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/valid" && r.Header.Get("Authorization") == "bearer validBearer" {
			resp, _ := json.Marshal(Response{Count: 1, Value: "validJWTToken"})
			_, _ = w.Write(resp)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}))
}

func TestFetchToken(t *testing.T) {
	testCases := []struct {
		name      string
		tokenURL  string
		bearer    string
		audience  string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid token",
			tokenURL:  "/valid",
			bearer:    "validBearer",
			audience:  "validAudience",
			wantToken: "validJWTToken",
			wantErr:   false,
		},
		{
			name:      "invalid token url",
			tokenURL:  "/invalid",
			bearer:    "validBearer",
			audience:  "validAudience",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "invalid bearer",
			tokenURL:  "/valid",
			bearer:    "invalidBearer",
			audience:  "validAudience",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "invalid url",
			tokenURL:  "invalidURL",
			bearer:    "validBearer",
			audience:  "validAudience",
			wantToken: "",
			wantErr:   true,
		},
	}

	server := createMockServer()
	defer server.Close()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotToken, err := fetchToken(server.URL+testCase.tokenURL, testCase.bearer, testCase.audience)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantToken, gotToken)
			}
		})
	}
}

func TestSubjects(t *testing.T) {
	tokenServer := createMockServer()
	defer tokenServer.Close()
	attestor := &Attestor{
		aud:      "projecturl",
		jwksURL:  tokenServer.URL,
		tokenURL: os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
	}

	subjects := attestor.Subjects()
	assert.NotNil(t, subjects)
	assert.Equal(t, 2, len(subjects))

	expectedSubjects := []string{"pipelineurl:" + attestor.PipelineUrl, "projecturl:" + attestor.ProjectUrl}
	for _, expectedSubject := range expectedSubjects {
		_, ok := subjects[expectedSubject]
		assert.True(t, ok, "Expected subject not found: %s", expectedSubject)
	}
	m := attestor.BackRefs()
	assert.NotNil(t, m)
	assert.Equal(t, 1, len(m))
}
