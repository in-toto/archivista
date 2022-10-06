// Copyright 2022 The Witness Contributors
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

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLResponse[T any] struct {
	Data   T              `json:"data,omitempty"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type graphQLRequestBody[TVars any] struct {
	Query     string `json:"query"`
	Variables TVars  `json:"variables,omitempty"`
}

func GraphQlQuery[TRes any, TVars any](ctx context.Context, baseUrl, query string, vars TVars) (TRes, error) {
	var response TRes
	queryUrl, err := url.JoinPath(baseUrl, "query")
	if err != nil {
		return response, err
	}

	requestBody := graphQLRequestBody[TVars]{
		Query:     query,
		Variables: vars,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", queryUrl, bytes.NewReader(reqBody))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	gqlRes := graphQLResponse[TRes]{}
	if err := dec.Decode(&gqlRes); err != nil {
		return response, err
	}

	if len(gqlRes.Errors) > 0 {
		return response, fmt.Errorf("graph ql query failed: %v", gqlRes.Errors)
	}

	return gqlRes.Data, nil
}
