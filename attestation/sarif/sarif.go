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

package sarif

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/owenrumney/go-sarif/sarif"
	"github.com/testifysec/go-witness/attestation"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/log"
)

const (
	Name    = "sarif"
	Type    = "https://witness.dev/attestations/sarif/v0.1"
	RunType = attestation.PostProductRunType
)

// This is a hacky way to create a compile time error in case the attestor
// doesn't implement the expected interfaces.
var (
	_ attestation.Attestor = &Attestor{}

	mimeTypes = []string{"text/plain", "application/json"}
)

func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor {
		return New()
	})
}

type Attestor struct {
	sarif.Report    `json:"report"`
	ReportFile      string               `json:"reportFileName"`
	ReportDigestSet cryptoutil.DigestSet `json:"reportDigestSet"`
}

func New() *Attestor {
	return &Attestor{}
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
	if err := a.getCandidate(ctx); err != nil {
		log.Debugf("(attestation/sarif) error getting candidate: %v", err)
		return err
	}

	return nil
}

func (a *Attestor) getCandidate(ctx *attestation.AttestationContext) error {
	products := ctx.Products()

	if len(products) == 0 {
		return fmt.Errorf("no products to attest")
	}

	for path, product := range products {
		for _, mimeType := range mimeTypes {
			if !strings.Contains(mimeType, product.MimeType) {
				continue
			}
		}

		newDigestSet, err := cryptoutil.CalculateDigestSetFromFile(path, ctx.Hashes())
		if newDigestSet == nil || err != nil {
			return fmt.Errorf("error calculating digest set from file: %s", path)
		}

		if !newDigestSet.Equal(product.Digest) {
			return fmt.Errorf("integrity error: product digest set does not match candidate digest set")
		}

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file: %s", path)
		}

		reportBytes, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("error reading file: %s", path)
		}

		//check to see if we can unmarshal into sarif type
		if err := json.Unmarshal(reportBytes, &a.Report); err != nil {
			log.Debugf("(attestation/sarif) error unmarshaling report: %v", err)
			continue
		}

		a.ReportFile = path
		a.ReportDigestSet = product.Digest

		return nil
	}
	return fmt.Errorf("no sarif file found")
}
