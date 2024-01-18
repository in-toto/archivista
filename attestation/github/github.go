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

	"github.com/davecgh/go-spew/spew"
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

func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor {
		return New()
	})
}

type ErrNotGitlab struct{}

func (e ErrNotGitlab) Error() string {
	return "not in a github ci job"
}

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

func New() *Attestor {
	return &Attestor{
		aud:      tokenAudience,
		jwksURL:  jwksURL,
		tokenURL: os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
	}
}

func (a *Attestor) Name() string {
	return Name
}

func (a *Attestor) Type() string {
	return Type
}

func (a *Attestor) RunType() attestation.RunType {
	return RunType
}

func (a *Attestor) Attest(ctx *attestation.AttestationContext) error {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return ErrNotGitlab{}
	}

	jwtString, err := fetchToken(a.tokenURL, os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"), "witness")
	if err != nil {
		return err
	}

	spew.Dump(jwtString)

	if jwtString != "" {
		a.JWT = jwt.New(jwt.WithToken(jwtString), jwt.WithJWKSUrl(a.jwksURL))
		if err := a.JWT.Attest(ctx); err != nil {
			return err
		}
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

func (a *Attestor) Subjects() map[string]cryptoutil.DigestSet {
	subjects := make(map[string]cryptoutil.DigestSet)
	hashes := []crypto.Hash{crypto.SHA256}
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

func fetchToken(tokenURL string, bearer string, audience string) (string, error) {
	client := &http.Client{}

	//add audient "&audience=witness" to the end of the tokenURL, parse it, and then add it to the query
	u, err := url.Parse(tokenURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
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
	body, err := readResponseBody(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse GithubTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.Value, nil
}

type GithubTokenResponse struct {
	Count int    `json:"count"`
	Value string `json:"value"`
}

func readResponseBody(body io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
