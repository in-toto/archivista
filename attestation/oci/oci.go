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

package oci

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
)

const (
	Name    = "oci"
	Type    = "https://witness.dev/attestations/oci/v0.1"
	RunType = attestation.PostProductRunType

	mimeTypes = "application/x-tar"
)

// This is a hacky way to create a compile time error in case the attestor
// doesn't implement the expected interfaces.
var (
	_ attestation.Attestor  = &Attestor{}
	_ attestation.Subjecter = &Attestor{}
)

func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor {
		return New()
	})
}

type Attestor struct {
	TarDigest      cryptoutil.DigestSet   `json:"tardigest"`
	Manifest       []Manifest             `json:"manifest"`
	ImageTags      []string               `json:"imagetags"`
	LayerDiffIDs   []cryptoutil.DigestSet `json:"diffids"`
	ImageID        cryptoutil.DigestSet   `json:"imageid"`
	ManifestRaw    []byte                 `json:"manifestraw"`
	ManifestDigest cryptoutil.DigestSet   `json:"manifestdigest"`
	tarFilePath    string                 `json:"-"`
}

type Manifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

func (m *Manifest) getImageID(ctx *attestation.AttestationContext, tarFilePath string) (cryptoutil.DigestSet, error) {
	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	tarReader := tar.NewReader(tarFile)
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.FileInfo().IsDir() {
			continue
		}

		if h.Name == m.Config {

			b := make([]byte, h.Size)
			_, err := tarReader.Read(b)
			if err != nil && err != io.EOF {
				return nil, err
			}

			imageID, err := cryptoutil.CalculateDigestSetFromBytes(b, ctx.Hashes())
			if err != nil {
				log.Debugf("(attestation/oci) error calculating image id: %w", err)
				return nil, err
			}

			return imageID, nil
		}
	}
	return nil, fmt.Errorf("could not find config in tar file")
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
		log.Debugf("(attestation/oci) error getting candidate: %w", err)
		return err
	}

	if err := a.parseMaifest(ctx); err != nil {
		log.Debugf("(attestation/oci) error parsing manifest: %w", err)
		return err
	}

	imageID, err := a.Manifest[0].getImageID(ctx, a.tarFilePath)
	if err != nil {
		log.Debugf("(attestation/oci) error getting image id: %w", err)
		return err
	}

	layerDiffIDs, err := a.Manifest[0].getLayerDIFFIDs(ctx, a.tarFilePath)
	if err != nil {
		return err
	}

	a.ImageID = imageID
	a.LayerDiffIDs = layerDiffIDs
	a.ImageTags = a.Manifest[0].RepoTags

	return nil
}

func (a *Attestor) getCandidate(ctx *attestation.AttestationContext) error {
	products := ctx.Products()

	if len(products) == 0 {
		return fmt.Errorf("no products to attest")
	}

	for path, product := range products {
		if !strings.Contains(mimeTypes, product.MimeType) {
			continue
		}

		newDigestSet, err := cryptoutil.CalculateDigestSetFromFile(path, ctx.Hashes())
		if newDigestSet == nil || err != nil {
			return fmt.Errorf("error calculating digest set from file: %s", path)
		}

		if !newDigestSet.Equal(product.Digest) {
			return fmt.Errorf("integrity error: product digest set does not match candidate digest set")
		}

		a.TarDigest = product.Digest

		a.tarFilePath = path
		return nil
	}
	return fmt.Errorf("no tar file found")
}

func (a *Attestor) parseMaifest(ctx *attestation.AttestationContext) error {

	f, err := os.Open(a.tarFilePath)
	if err != nil {
		err = fmt.Errorf("error opening tar file: %w", err)
		return err
	}

	tarReader := tar.NewReader(f)
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if h.FileInfo().IsDir() {
			continue
		}
		if h.Name == "manifest.json" {
			a.ManifestRaw = make([]byte, h.Size)
			_, err = tarReader.Read(a.ManifestRaw)
			if err != nil || err == io.EOF {
				break
			}
			break
		}
	}

	manifestDigest, err := cryptoutil.CalculateDigestSetFromBytes(a.ManifestRaw, ctx.Hashes())
	if err != nil {
		return err
	}

	a.ManifestDigest = manifestDigest

	err = json.Unmarshal(a.ManifestRaw, &a.Manifest)
	if err != nil {
		return err
	}

	return nil
}

func (a *Attestor) Subjects() map[string]cryptoutil.DigestSet {
	hashes := []cryptoutil.DigestValue{{Hash: crypto.SHA256}}
	subj := make(map[string]cryptoutil.DigestSet)
	subj[fmt.Sprintf("manifestdigest:%s", a.ManifestDigest[cryptoutil.DigestValue{Hash: crypto.SHA256}])] = a.ManifestDigest
	subj[fmt.Sprintf("tardigest:%s", a.TarDigest[cryptoutil.DigestValue{Hash: crypto.SHA256}])] = a.TarDigest
	subj[fmt.Sprintf("imageid:%s", a.ImageID[cryptoutil.DigestValue{Hash: crypto.SHA256}])] = a.ImageID

	//image tags
	for _, tag := range a.ImageTags {
		hash, err := cryptoutil.CalculateDigestSetFromBytes([]byte(tag), hashes)
		if err != nil {
			log.Debugf("(attestation/oci) error calculating image tag: %w", err)
			continue
		}
		subj[fmt.Sprintf("imagetag:%s", tag)] = hash
	}

	//diff ids
	for layer := range a.LayerDiffIDs {
		subj[fmt.Sprintf("layerdiffid%02d:%s", layer, a.LayerDiffIDs[layer][cryptoutil.DigestValue{Hash: crypto.SHA256}])] = a.LayerDiffIDs[layer]
	}
	return subj
}

func (m *Manifest) getLayerDIFFIDs(ctx *attestation.AttestationContext, tarFilePath string) ([]cryptoutil.DigestSet, error) {
	var layerDiffIDs []cryptoutil.DigestSet

	tarFile, err := os.Open(tarFilePath)
	if err != nil {
		return nil, err
	}
	defer tarFile.Close()

	tarReader := tar.NewReader(tarFile)
	for {
		h, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.FileInfo().IsDir() {
			continue
		}
		for _, layerFile := range m.Layers {
			if h.Name == layerFile {
				b := make([]byte, h.Size)

				_, err := tarReader.Read(b)
				if err != nil && err != io.EOF {
					return nil, err
				}

				contentType := http.DetectContentType(b)
				if contentType == "application/x-gzip" {
					breader, err := gzip.NewReader(bytes.NewReader(b))
					if err != nil {
						return nil, err
					}
					defer breader.Close()
					c, err := io.ReadAll(breader)
					if err != nil {
						return nil, err
					}
					layerDiffID, err := cryptoutil.CalculateDigestSetFromBytes(c, ctx.Hashes())
					if err != nil {
						return nil, err
					}
					layerDiffIDs = append(layerDiffIDs, layerDiffID)

				} else {
					layerDiffID, err := cryptoutil.CalculateDigestSetFromBytes(b, ctx.Hashes())
					if err != nil {
						return nil, err
					}
					layerDiffIDs = append(layerDiffIDs, layerDiffID)
				}

			}

		}
	}
	return layerDiffIDs, nil
}
