// Copyright 2023 The Witness Contributors
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

package fulcio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchToken(t *testing.T) {
	tokenURL := "https://example.com/token"
	bearer := "some-bearer-token"
	audience := "witness"

	// Test empty input.
	_, err := fetchToken("", bearer, audience)
	require.Error(t, err)

	_, err = fetchToken(tokenURL, "", audience)
	require.Error(t, err)

	_, err = fetchToken(tokenURL, bearer, "")
	require.Error(t, err)

	// Test invalid input.
	u, _ := url.Parse(tokenURL)
	q := u.Query()
	q.Set("audience", "other-audience")
	u.RawQuery = q.Encode()

	_, err = fetchToken(u.String(), bearer, audience)
	require.Error(t, err)

	// Test valid input.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "bearer "+bearer {
			http.Error(w, "Invalid bearer token", http.StatusUnauthorized)
			return
		}
		if r.URL.Query().Get("audience") != audience {
			http.Error(w, "Invalid audience", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `{"count": 1, "value": "some-token"}`)
	}))
	defer server.Close()

	tokenURL = server.URL + "/token"

	token, err := fetchToken(tokenURL, bearer, audience)
	require.NoError(t, err)
	require.Equal(t, "some-token", token)
}
