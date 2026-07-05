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

package rstuf

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/in-toto/archivista/pkg/config"
	"github.com/sirupsen/logrus"
)

type RSTUF struct {
	Host string
}

type Publisher interface {
	Publish(ctx context.Context, gitoid string, payload []byte) error
}

func (r *RSTUF) parseRSTUFPayload(gitoid string, payload []byte) ([]byte, error) {
	objHash := sha256.Sum256(payload)
	// custom := make(map[string]any)
	// custom["gitoid"] = gitoid
	artifacts := []Artifact{
		{
			Path: gitoid,
			Info: ArtifactInfo{
				Length: len(payload),
				Hashes: Hashes{
					Sha256: hex.EncodeToString(objHash[:]),
				},
				// Custom: custom,
			},
		},
	}

	artifactPayload := ArtifactPayload{
		Artifacts:         artifacts,
		AddTaskIDToCustom: false,
		PublishTargets:    true,
	}

	payloadBytes, err := json.Marshal(artifactPayload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}
	return payloadBytes, nil
}

func (r *RSTUF) Publish(ctx context.Context, gitoid string, payload []byte) error {
	// this publisher allows integration with the RSTUF project to store
	// the attestation and policy in the TUF metadata.
	// this TUF metadata can be used to build truste when distributing the
	// attestations and policies.
	// Convert payload to JSON
	url := r.Host + "/api/v1/artifacts"

	payloadBytes, err := r.parseRSTUFPayload(gitoid, payload)
	if err != nil {
		return fmt.Errorf("error parsing payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Add any additional headers or authentication if needed

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		logb, _ := httputil.DumpResponse(resp, true)
		logrus.Errorf("error body from RSTUF: %v", string(logb))
		return fmt.Errorf("error response from RSTUF: %v", err)
	}

	// Handle the response as needed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error reading response body: %v", err)
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Errorf("error unmarshaling response: %v", err)
	}
	logrus.Debugf("RSTUF task id: %v", response.Data.TaskId)
	// TODO: monitor RSTUF task id for completion
	return nil
}

func NewPublisher(config *config.Config) Publisher {
	return &RSTUF{
		Host: config.PublisherRstufHost,
	}
}
