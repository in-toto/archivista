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

package attestation

import (
	"fmt"

	"github.com/testifysec/go-witness/cryptoutil"
)

var (
	attestationsByName = map[string]AttestorRegistration{}
	attestationsByType = map[string]AttestorRegistration{}
	attestationsByRun  = map[RunType]AttestorRegistration{}
)

type Configurer interface {
	Description() string
	Name() string
}

type Configurable interface {
	int | string | []string
}

type ConfigOption[T Configurable] struct {
	name        string
	description string
	defaultVal  T
	setter      func(Attestor, T) (Attestor, error)
}

func (co ConfigOption[T]) Name() string {
	return co.name
}

func (co ConfigOption[T]) DefaultVal() T {
	return co.defaultVal
}

func (co ConfigOption[T]) Description() string {
	return co.description
}

func (co ConfigOption[T]) Setter() func(Attestor, T) (Attestor, error) {
	return co.setter
}

func IntConfigOption(name, description string, defaultVal int, setter func(Attestor, int) (Attestor, error)) ConfigOption[int] {
	return ConfigOption[int]{
		name,
		description,
		defaultVal,
		setter,
	}
}

func StringConfigOption(name, description string, defaultVal string, setter func(Attestor, string) (Attestor, error)) ConfigOption[string] {
	return ConfigOption[string]{
		name,
		description,
		defaultVal,
		setter,
	}
}

func StringSliceConfigOption(name, description string, defaultVal []string, setter func(Attestor, []string) (Attestor, error)) ConfigOption[[]string] {
	return ConfigOption[[]string]{
		name,
		description,
		defaultVal,
		setter,
	}
}

type AttestorRegistration struct {
	Factory AttestorFactory
	Name    string
	Type    string
	RunType RunType
	Options []Configurer
}

type Attestor interface {
	Name() string
	Type() string
	RunType() RunType
	Attest(ctx *AttestationContext) error
}

// Subjecter allows attestors to expose bits of information that will be added to
// the in-toto statement as subjects. External services such as Rekor and Archivist
// use in-toto subjects as indexes back to attestations.
type Subjecter interface {
	Subjects() map[string]cryptoutil.DigestSet
}

// Materialer allows attestors to communicate about materials that were observed
// while the attestor executed. For example the material attestor records the hashes
// of all files before a command is run.
type Materialer interface {
	Materials() map[string]cryptoutil.DigestSet
}

// Producer allows attestors to communicate that some product was created while the
// attestor executed. For example the product attestor runs after a command run and
// finds files that did not exist in the working directory prior to the command's
// execution.
type Producer interface {
	Products() map[string]Product
}

// BackReffer allows attestors to indicate which of their subjects are good candidates
// to find related attestations.  For example the git attestor's commit hash subject
// is a good candidate to find all attestation collections that also refer to a specific
// git commit.
type BackReffer interface {
	BackRefs() map[string]cryptoutil.DigestSet
}

type AttestorFactory func() Attestor

type ErrAttestationNotFound string

func (e ErrAttestationNotFound) Error() string {
	return fmt.Sprintf("attestation not found: %v", string(e))
}

func RegisterAttestation(name, predicateType string, run RunType, factoryFunc AttestorFactory, opts ...Configurer) {
	registrationEntry := AttestorRegistration{
		Name:    name,
		Type:    predicateType,
		Factory: factoryFunc,
		RunType: run,
		Options: opts,
	}

	attestationsByName[name] = registrationEntry
	attestationsByType[predicateType] = registrationEntry
	attestationsByRun[run] = registrationEntry
}

func FactoryByType(uri string) (AttestorFactory, bool) {
	registrationEntry, ok := attestationsByType[uri]
	return registrationEntry.Factory, ok
}

func FactoryByName(name string) (AttestorFactory, bool) {
	registrationEntry, ok := attestationsByName[name]
	return registrationEntry.Factory, ok
}

func Attestors(nameOrTypes []string) ([]Attestor, error) {
	attestors := make([]Attestor, 0)
	for _, nameOrType := range nameOrTypes {
		factory, ok := FactoryByName(nameOrType)
		if !ok {
			factory, ok = FactoryByType(nameOrType)
			if !ok {
				return nil, ErrAttestationNotFound(nameOrType)
			}
		}

		attestor := factory()
		opts := AttestorOptions(nameOrType)
		attestor, err := setDefaultVals(attestor, opts)
		if err != nil {
			return nil, err
		}

		attestors = append(attestors, attestor)
	}

	return attestors, nil
}

func AttestorOptions(nameOrType string) []Configurer {
	entry, ok := attestationsByName[nameOrType]
	if !ok {
		entry = attestationsByType[nameOrType]
	}

	return entry.Options
}

func RegistrationEntries() []AttestorRegistration {
	results := make([]AttestorRegistration, 0, len(attestationsByName))
	for _, registration := range attestationsByName {
		results = append(results, registration)
	}

	return results
}

func setDefaultVals(attestor Attestor, opts []Configurer) (Attestor, error) {
	var err error

	for _, opt := range opts {
		switch o := opt.(type) {
		case ConfigOption[int]:
			attestor, err = o.Setter()(attestor, o.DefaultVal())
		case ConfigOption[string]:
			attestor, err = o.Setter()(attestor, o.DefaultVal())
		case ConfigOption[[]string]:
			attestor, err = o.Setter()(attestor, o.DefaultVal())
		}

		if err != nil {
			return attestor, err
		}
	}

	return attestor, nil
}
