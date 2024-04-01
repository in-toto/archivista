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

package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// see https://developer.hashicorp.com/vault/api-docs/secret/pki#issuing-certificates
// for information on the following structs and requests

type issueRequest struct {
	CommonName           string        `json:"common_name,omitempty"`
	AltNames             []string      `json:"alt_names,omitempty"`
	Ttl                  time.Duration `json:"ttl,omitempty"`
	RemoveRootsFromChain bool          `json:"remove_roots_from_chain,omitempty"`
}

type issueResponseData struct {
	Certificate    string   `json:"certificate"`
	IssuingCa      string   `json:"issuing_ca"`
	CaChain        []string `json:"ca_chain"`
	PrivateKey     string   `json:"private_key"`
	PrivateKeyType string   `json:"private_key_type"`
	SerialNumber   string   `json:"serial_number"`
}

type issueResponse struct {
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration int               `json:"lease_duration"`
	Warnings      []string          `json:"warnings"`
	Data          issueResponseData `json:"data"`
}

func (vsp *VaultSignerProvider) requestCertificate(ctx context.Context) (issueResponse, error) {
	url, err := url.JoinPath(vsp.url, "v1", vsp.pkiSecretsEnginePath, "issue", vsp.role)
	if err != nil {
		return issueResponse{}, err
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(issueRequest{
		CommonName:           vsp.commonName,
		AltNames:             vsp.altNames,
		Ttl:                  vsp.ttl,
		RemoveRootsFromChain: true,
	}); err != nil {
		return issueResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return issueResponse{}, err
	}

	req.Header.Set("X-Vault-Token", vsp.token)
	if len(vsp.namespace) > 0 {
		req.Header.Set("X-Vault-Namespace", vsp.namespace)
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return issueResponse{}, err
	}

	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return issueResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return issueResponse{}, fmt.Errorf("failed to issue new certificate: %s", respBytes)
	}

	issueResp := issueResponse{}
	if err := json.Unmarshal(respBytes, &issueResp); err != nil {
		return issueResp, err
	}

	return issueResp, nil
}
