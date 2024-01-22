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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/in-toto/go-witness/dsse"
)

type StoreResponse struct {
	Gitoid string `json:"gitoid"`
}

func Store(ctx context.Context, baseURL string, envelope dsse.Envelope) (StoreResponse, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(envelope); err != nil {
		return StoreResponse{}, err
	}

	return StoreWithReader(ctx, baseURL, buf)
}

func StoreWithReader(ctx context.Context, baseURL string, r io.Reader) (StoreResponse, error) {
	return StoreWithReaderWithHTTPClient(ctx, &http.Client{}, baseURL, r)
}

func StoreWithReaderWithHTTPClient(ctx context.Context, client *http.Client, baseURL string, r io.Reader) (StoreResponse, error) {
	uploadPath, err := url.JoinPath(baseURL, "upload")
	if err != nil {
		return StoreResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uploadPath, r)
	if err != nil {
		return StoreResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return StoreResponse{}, err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return StoreResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return StoreResponse{}, errors.New(string(bodyBytes))
	}

	storeResp := StoreResponse{}
	if err := json.Unmarshal(bodyBytes, &storeResp); err != nil {
		return StoreResponse{}, err
	}

	return storeResp, nil
}
