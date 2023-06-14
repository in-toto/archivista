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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func fetchToken(tokenURL string, bearer string, audience string) (string, error) {
	if tokenURL == "" || bearer == "" || audience == "" {
		return "", fmt.Errorf("tokenURL, bearer, and audience cannot be empty")
	}

	client := &http.Client{}

	//add audient "&audience=witness" to the end of the tokenURL, parse it, and then add it to the query
	u, err := url.Parse(tokenURL)
	if err != nil {
		return "", err
	}

	//check to see if the tokenURL already has a query with an audience
	//if it does throw an error
	q := u.Query()
	if q.Get("audience") != audience && q.Get("audience") != "" {
		return "", fmt.Errorf("api error: tokenURL already has an audience, %s, and it does not match the audience, %s", q.Get("audience"), audience)
	}

	q.Add("audience", audience)
	u.RawQuery = q.Encode()

	reqURL := u.String()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "bearer "+bearer)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse struct {
		Count int    `json:"count"`
		Value string `json:"value"`
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Value, nil
}
