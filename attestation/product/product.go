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

package product

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobwas/glob"
	"github.com/testifysec/go-witness/attestation"
	"github.com/testifysec/go-witness/attestation/file"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/registry"
)

const (
	Name    = "product"
	Type    = "https://witness.dev/attestations/product/v0.1"
	RunType = attestation.ProductRunType

	defaultIncludeGlob = "*"
	defaultExcludeGlob = ""
)

// This is a hacky way to create a compile time error in case the attestor
// doesn't implement the expected interfaces.
var (
	_ attestation.Attestor  = &Attestor{}
	_ attestation.Subjecter = &Attestor{}
	_ attestation.Producer  = &Attestor{}
)

func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor { return New() },
		registry.StringConfigOption(
			"include-glob",
			"Pattern to use when recording products. Files that match this pattern will be included as subjects on the attestation.",
			defaultIncludeGlob,
			func(a attestation.Attestor, includeGlob string) (attestation.Attestor, error) {
				prodAttestor, ok := a.(*Attestor)
				if !ok {
					return a, fmt.Errorf("unexpected attestor type: %T is not a product attestor", a)
				}

				WithIncludeGlob(includeGlob)(prodAttestor)
				return prodAttestor, nil
			},
		),
		registry.StringConfigOption(
			"exclude-glob",
			"Pattern to use when recording products. Files that match this pattern will be excluded as subjects on the attestation.",
			defaultExcludeGlob,
			func(a attestation.Attestor, excludeGlob string) (attestation.Attestor, error) {
				prodAttestor, ok := a.(*Attestor)
				if !ok {
					return a, fmt.Errorf("unexpected attestor type: %T is not a product attestor", a)
				}

				WithExcludeGlob(excludeGlob)(prodAttestor)
				return prodAttestor, nil
			},
		),
	)
}

type Option func(*Attestor)

func WithIncludeGlob(glob string) Option {
	return func(a *Attestor) {
		a.includeGlob = glob
	}
}

func WithExcludeGlob(glob string) Option {
	return func(a *Attestor) {
		a.excludeGlob = glob
	}
}

type Attestor struct {
	products            map[string]attestation.Product
	baseArtifacts       map[string]cryptoutil.DigestSet
	includeGlob         string
	compiledIncludeGlob glob.Glob
	excludeGlob         string
	compiledExcludeGlob glob.Glob
}

func fromDigestMap(digestMap map[string]cryptoutil.DigestSet) map[string]attestation.Product {
	products := make(map[string]attestation.Product)
	for fileName, digestSet := range digestMap {
		mimeType := "unknown"
		f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
		if err == nil {
			mimeType, err = getFileContentType(f)
			if err != nil {
				mimeType = "unknown"
			}
			f.Close()
		}

		products[fileName] = attestation.Product{
			MimeType: mimeType,
			Digest:   digestSet,
		}
	}

	return products
}

func (a Attestor) Name() string {
	return Name
}

func (a Attestor) Type() string {
	return Type
}

func (rc *Attestor) RunType() attestation.RunType {
	return RunType
}

func New(opts ...Option) *Attestor {
	a := &Attestor{
		includeGlob: defaultIncludeGlob,
		excludeGlob: defaultExcludeGlob,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Attestor) Attest(ctx *attestation.AttestationContext) error {
	compiledIncludeGlob, err := glob.Compile(a.includeGlob)
	if err != nil {
		return err
	}
	a.compiledIncludeGlob = compiledIncludeGlob

	compiledExcludeGlob, err := glob.Compile(a.excludeGlob)
	if err != nil {
		return err
	}
	a.compiledExcludeGlob = compiledExcludeGlob

	a.baseArtifacts = ctx.Materials()
	products, err := file.RecordArtifacts(ctx.WorkingDir(), a.baseArtifacts, ctx.Hashes(), map[string]struct{}{})
	if err != nil {
		return err
	}

	a.products = fromDigestMap(products)
	return nil
}

func (a *Attestor) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.products)
}

func (a *Attestor) UnmarshalJSON(data []byte) error {
	prods := make(map[string]attestation.Product)
	if err := json.Unmarshal(data, &prods); err != nil {
		return err
	}

	a.products = prods
	return nil
}

func (a *Attestor) Products() map[string]attestation.Product {
	return a.products
}

func (a *Attestor) Subjects() map[string]cryptoutil.DigestSet {
	subjects := make(map[string]cryptoutil.DigestSet)
	for productName, product := range a.products {
		if a.compiledExcludeGlob != nil && a.compiledExcludeGlob.Match(productName) {
			continue
		}

		if a.compiledIncludeGlob != nil && !a.compiledIncludeGlob.Match(productName) {
			continue
		}

		subjects[fmt.Sprintf("file:%v", productName)] = product.Digest
	}

	return subjects
}

func getFileContentType(file *os.File) (string, error) {
	// Read up to 512 bytes from the file.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Try to detect the content type using http.DetectContentType().
	contentType := http.DetectContentType(buffer)

	// If the content type is application/octet-stream, try to detect the content type using a file signature.
	if contentType == "application/octet-stream" {
		// Try to match the file signature to a content type.
		if signature, _ := getFileSignature(buffer); signature != "application/octet-stream" {
			contentType = signature
		} else if extension := filepath.Ext(file.Name()); extension != "" {
			contentType = mime.TypeByExtension(extension)
		}
	}

	return contentType, nil
}

// getFileSignature tries to match the file signature to a content type.
func getFileSignature(buffer []byte) (string, error) {
	// Create a new buffer with a length of 512 bytes and copy the data from the input buffer into the new buffer to prevent out of bounds errors.
	newBuffer := make([]byte, 512)
	copy(newBuffer, buffer)

	var signature string
	switch {
	// https://en.wikipedia.org/wiki/List_of_file_signatures
	case buffer[257] == 0x75 && buffer[258] == 0x73 && buffer[259] == 0x74 && buffer[260] == 0x61 && buffer[261] == 0x72:
		signature = "application/x-tar"
	case buffer[0] == 0x25 && buffer[1] == 0x50 && buffer[2] == 0x44 && buffer[3] == 0x46 && buffer[4] == 0x2D:
		signature = "application/pdf"
	default:
		// If the file signature is not recognized, return application/octet-stream by default
		signature = "application/octet-stream"
	}

	return signature, nil
}
