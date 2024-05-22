// Copyright 2024 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package artifactstore

import (
	"crypto"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/in-toto/go-witness/cryptoutil"
	"gopkg.in/yaml.v3"
)

// Config represents the available artifacts within the store
type Config struct {
	Artifacts map[string]Artifact `json:"artifacts"`
}

// Artifact represents an artifact and it's available versions in the store
type Artifact struct {
	Versions map[string]Version `json:"versions"`
}

// Version represents a version of an Artifact (ex v0.2.0) with the available Distributions of version
type Version struct {
	Distributions map[string]Distribution `json:"distributions"`
	Description   string                  `json:"description"`
}

// Distribution is a specific distribution of a Version of an Artifact(ex linux-amd64)
type Distribution struct {
	FileLocation string `json:"-"`
	SHA256Digest string `json:"sha256digest"`
}

// Store is an artifact store served from Archivista
type Store struct {
	config Config
}

type Option func(*Store) error

// WithConfig creates a Store with the provided config
func WithConfig(config Config) Option {
	return func(as *Store) error {
		as.config = config
		return nil
	}
}

// WithConfigFile creates a Store with a config loaded from a yaml file on disk
func WithConfigFile(configPath string) Option {
	return func(as *Store) error {
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}

		config := Config{}
		if err := yaml.Unmarshal(configBytes, &config); err != nil {
			return err
		}

		return WithConfig(config)(as)
	}
}

// New creates a new Store with the provided options
func New(opts ...Option) (Store, error) {
	as := Store{}
	errs := make([]error, 0)
	for _, opt := range opts {
		if err := opt(&as); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return Store{}, errors.Join(errs...)
	}

	if err := verifyConfig(as); err != nil {
		return Store{}, err
	}

	return as, nil
}

// verifyConfig ensures that each file exists on disk and that the sha256sum of the
// files on disk match those of the config
func verifyConfig(as Store) error {
	errs := make([]error, 0)
	for artifactName, artifact := range as.config.Artifacts {
		for versionString, version := range artifact.Versions {
			for distroString, distro := range version.Distributions {
				if _, err := os.Stat(distro.FileLocation); err != nil {
					errs = append(errs, fmt.Errorf("%v version %v-%v does not exist on disk: %w", artifactName, versionString, distroString, err))
					continue
				}

				digestSet, err := cryptoutil.CalculateDigestSetFromFile(distro.FileLocation, []cryptoutil.DigestValue{{Hash: crypto.SHA256}})
				if err != nil {
					errs = append(errs, fmt.Errorf("could not calculate sha256 digest for %v version %v-%v: %w", artifactName, versionString, distroString, err))
				}

				sha256Digest := digestSet[cryptoutil.DigestValue{Hash: crypto.SHA256, GitOID: false}]
				if !strings.EqualFold(sha256Digest, distro.SHA256Digest) {
					errs = append(errs, fmt.Errorf("sha256 digest of %v version %v-%v does not match config: got %v, expected %v", artifactName, versionString, distroString, sha256Digest, distro.SHA256Digest))
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Artifacts returns a copy of all the Store's Artifacts
func (as Store) Artifacts() map[string]Artifact {
	out := make(map[string]Artifact)
	for artifactString, artifact := range as.config.Artifacts {
		out[artifactString] = artifact
	}

	return out
}

// Versions returns all of the available Versions for an Artifact.
func (as Store) Versions(artifact string) (map[string]Version, bool) {
	out := make(map[string]Version)
	a, ok := as.config.Artifacts[artifact]
	if !ok {
		return out, false
	}

	for verString, version := range a.Versions {
		out[verString] = version
	}

	return out, ok
}

// Version returns a specific Version for an artifact, if it exists
func (as Store) Version(artifact, version string) (Version, bool) {
	a, ok := as.config.Artifacts[artifact]
	if !ok {
		return Version{}, false
	}

	v, ok := a.Versions[version]
	return v, ok
}

// Distributions returns all of the available Distributions for a specified Version of an Artifact
func (as Store) Distributions(artifact, version string) (map[string]Distribution, bool) {
	out := make(map[string]Distribution)
	a, ok := as.config.Artifacts[artifact]
	if !ok {
		return out, false
	}

	vers, ok := a.Versions[version]
	if !ok {
		return out, ok
	}

	for distroString, distro := range vers.Distributions {
		out[distroString] = distro
	}

	return out, true
}

// Distribution returns the entry for a specific distribution for a specific version
func (as Store) Distribution(artifact, version, distribution string) (Distribution, bool) {
	a, ok := as.config.Artifacts[artifact]
	if !ok {
		return Distribution{}, false
	}

	vers, ok := a.Versions[version]
	if !ok {
		return Distribution{}, false
	}

	distro, ok := vers.Distributions[distribution]
	return distro, ok
}
