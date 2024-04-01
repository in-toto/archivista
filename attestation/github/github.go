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
	"bytes"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/attestation/jwt"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
)

const (
	Name    = "github"
	Type    = "https://witness.dev/attestations/github/v0.1"
	RunType = attestation.PreMaterialRunType
)

const (
	tokenAudience = "witness"
	jwksURL       = "https://token.actions.githubusercontent.com/.well-known/jwks"
)

// This is a hacky way to create a compile time error in case the attestor
// doesn't implement the expected interfaces.
var (
	_ attestation.Attestor   = &Attestor{}
	_ attestation.Subjecter  = &Attestor{}
	_ attestation.BackReffer = &Attestor{}
)

// init registers the github attestor.
func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor {
		return New()
	})
}

// ErrNotGitlab is an error type that indicates the environment is not a github ci job.
type ErrNotGitlab struct{}

// Error returns the error message for ErrNotGitlab.
func (e ErrNotGitlab) Error() string {
	return "not in a github ci job"
}

// Attestor is a struct that holds the necessary information for github attestation.
type Attestor struct {
	JWT          *jwt.Attestor `json:"jwt,omitempty"`
	CIConfigPath string        `json:"ciconfigpath"`
	PipelineID   string        `json:"pipelineid"`
	PipelineName string        `json:"pipelinename"`
	PipelineUrl  string        `json:"pipelineurl"`
	ProjectUrl   string        `json:"projecturl"`
	RunnerID     string        `json:"runnerid"`
	CIHost       string        `json:"cihost"`
	CIServerUrl  string        `json:"ciserverurl"`
	RunnerArch   string        `json:"runnerarch"`
	RunnerOS     string        `json:"runneros"`

	jwksURL  string
	tokenURL string
	aud      string
}

// New creates and returns a new github attestor.
func New() *Attestor {
	return &Attestor{
		aud:      tokenAudience,
		jwksURL:  jwksURL,
		tokenURL: os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
	}
}

// Name returns the name of the attestor.
func (a *Attestor) Name() string {
	return Name
}

// Type returns the type of the attestor.
func (a *Attestor) Type() string {
	return Type
}

// RunType returns the run type of the attestor.
func (a *Attestor) RunType() attestation.RunType {
	return RunType
}

// Attest performs the attestation for the github environment.
func (a *Attestor) Attest(ctx *attestation.AttestationContext) error {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return ErrNotGitlab{}
	}

	jwtString, err := fetchToken(a.tokenURL, os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"), "witness")
	if err != nil {
		return fmt.Errorf("error on fetching token %w", err)
	}

	if jwtString == "" {
		return fmt.Errorf("empty JWT string")
	}

	a.JWT = jwt.New(jwt.WithToken(jwtString), jwt.WithJWKSUrl(a.jwksURL))
	if err := a.JWT.Attest(ctx); err != nil {
		return fmt.Errorf("failed to attest github jwt: %w", err)
	}

	a.CIServerUrl = os.Getenv("GITHUB_SERVER_URL")
	a.CIConfigPath = os.Getenv("GITHUB_ACTION_PATH")

	a.PipelineID = os.Getenv("GITHUB_RUN_ID")
	a.PipelineName = os.Getenv("GITHUB_WORKFLOW")

	a.ProjectUrl = fmt.Sprintf("%s/%s", a.CIServerUrl, os.Getenv("GITHUB_REPOSITORY"))
	a.RunnerID = os.Getenv("RUNNER_NAME")
	a.RunnerArch = os.Getenv("RUNNER_ARCH")
	a.RunnerOS = os.Getenv("RUNNER_OS")
	a.PipelineUrl = fmt.Sprintf("%s/actions/runs/%s", a.ProjectUrl, a.PipelineID)
	return nil
}

// Subjects returns a map of subjects and their corresponding digest sets.
func (a *Attestor) Subjects() map[string]cryptoutil.DigestSet {
	subjects := make(map[string]cryptoutil.DigestSet)
	hashes := []cryptoutil.DigestValue{{Hash: crypto.SHA256}}
	if pipelineSubj, err := cryptoutil.CalculateDigestSetFromBytes([]byte(a.PipelineUrl), hashes); err == nil {
		subjects[fmt.Sprintf("pipelineurl:%v", a.PipelineUrl)] = pipelineSubj
	} else {
		log.Debugf("(attestation/github) failed to record github pipelineurl subject: %w", err)
	}

	if projectSubj, err := cryptoutil.CalculateDigestSetFromBytes([]byte(a.ProjectUrl), hashes); err == nil {
		subjects[fmt.Sprintf("projecturl:%v", a.ProjectUrl)] = projectSubj
	} else {
		log.Debugf("(attestation/github) failed to record github projecturl subject: %w", err)
	}

	return subjects
}

// BackRefs returns a map of back references and their corresponding digest sets.
func (a *Attestor) BackRefs() map[string]cryptoutil.DigestSet {
	backRefs := make(map[string]cryptoutil.DigestSet)
	for subj, ds := range a.Subjects() {
		if strings.HasPrefix(subj, "pipelineurl:") {
			backRefs[subj] = ds
			break
		}
	}

	return backRefs
}

// fetchToken fetches the token from the given URL.
func fetchToken(tokenURL string, bearer string, audience string) (string, error) {
	client := &http.Client{}

	// add audience "&audience=witness" to the end of the tokenURL, parse it, and then add it to the query
	u, err := url.Parse(tokenURL)
	if err != nil {
		return "", fmt.Errorf("error on parsing token url %w", err)
	}

	q := u.Query()
	q.Add("audience", audience)
	u.RawQuery = q.Encode()

	reqURL := u.String()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("error on creating request %w", err)
	}
	req.Header.Add("Authorization", "bearer "+bearer)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on request %w", err)
	}
	defer resp.Body.Close()
	body, err := readResponseBody(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error on reading response body %w", err)
	}

	var tokenResponse GithubTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("error on unmarshaling token response %w", err)
	}

	return tokenResponse.Value, nil
}

// GithubTokenResponse is a struct that holds the response from the github token request.
type GithubTokenResponse struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

// readResponseBody reads the response body and returns it as a byte slice.
func readResponseBody(body io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
