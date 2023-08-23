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

package attestation

import (
	"context"
	"crypto"
	"fmt"
	"os"
	"time"

	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/log"
)

type RunType string

const (
	PreMaterialRunType RunType = "prematerial"
	MaterialRunType    RunType = "material"
	ExecuteRunType     RunType = "execute"
	ProductRunType     RunType = "product"
	PostProductRunType RunType = "postproduct"
)

func (r RunType) String() string {
	return string(r)
}

type ErrInvalidOption struct {
	Option string
	Reason string
}

func (e ErrInvalidOption) Error() string {
	return fmt.Sprintf("invalid value for option %v: %v", e.Option, e.Reason)
}

type AttestationContextOption func(ctx *AttestationContext)

func WithContext(ctx context.Context) AttestationContextOption {
	return func(actx *AttestationContext) {
		actx.ctx = ctx
	}
}

func WithHashes(hashes []crypto.Hash) AttestationContextOption {
	return func(ctx *AttestationContext) {
		if len(hashes) > 0 {
			ctx.hashes = hashes
		}
	}
}

func WithWorkingDir(workingDir string) AttestationContextOption {
	return func(ctx *AttestationContext) {
		if workingDir != "" {
			ctx.workingDir = workingDir
		}
	}
}

type CompletedAttestor struct {
	Attestor  Attestor
	StartTime time.Time
	EndTime   time.Time
	Error     error
}

type AttestationContext struct {
	ctx                context.Context
	attestors          []Attestor
	workingDir         string
	hashes             []crypto.Hash
	completedAttestors []CompletedAttestor
	products           map[string]Product
	materials          map[string]cryptoutil.DigestSet
}

type Product struct {
	MimeType string               `json:"mime_type"`
	Digest   cryptoutil.DigestSet `json:"digest"`
}

func NewContext(attestors []Attestor, opts ...AttestationContextOption) (*AttestationContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ctx := &AttestationContext{
		ctx:        context.Background(),
		attestors:  attestors,
		workingDir: wd,
		hashes:     []crypto.Hash{crypto.SHA256},
		materials:  make(map[string]cryptoutil.DigestSet),
		products:   make(map[string]Product),
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx, nil
}

func (ctx *AttestationContext) RunAttestors() error {
	preAttestors := []Attestor{}
	materialAttestors := []Attestor{}
	exeucteAttestors := []Attestor{}
	productAttestors := []Attestor{}
	postAttestors := []Attestor{}

	for _, attestor := range ctx.attestors {
		switch attestor.RunType() {
		case PreMaterialRunType:
			preAttestors = append(preAttestors, attestor)

		case MaterialRunType:
			materialAttestors = append(materialAttestors, attestor)

		case ExecuteRunType:
			exeucteAttestors = append(exeucteAttestors, attestor)

		case ProductRunType:
			productAttestors = append(productAttestors, attestor)

		case PostProductRunType:
			postAttestors = append(postAttestors, attestor)

		default:
			return ErrInvalidOption{
				Option: "attestor.RunType",
				Reason: fmt.Sprintf("unknown run type %v", attestor.RunType()),
			}
		}
	}

	for _, attestor := range preAttestors {
		if err := ctx.runAttestor(attestor); err != nil {
			return err
		}
	}

	for _, attestor := range materialAttestors {
		if err := ctx.runAttestor(attestor); err != nil {
			return err
		}
	}

	for _, attestor := range exeucteAttestors {
		if err := ctx.runAttestor(attestor); err != nil {
			return err
		}
	}

	for _, attestor := range productAttestors {
		if err := ctx.runAttestor(attestor); err != nil {
			return err
		}
	}

	for _, attestor := range postAttestors {
		if err := ctx.runAttestor(attestor); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *AttestationContext) runAttestor(attestor Attestor) error {
	log.Infof("Starting %v attestor...", attestor.Name())
	startTime := time.Now()
	if err := attestor.Attest(ctx); err != nil {
		log.Errorf("Error running %v attestor: %v", attestor.Name(), err)
		ctx.completedAttestors = append(ctx.completedAttestors, CompletedAttestor{
			Attestor:  attestor,
			StartTime: startTime,
			EndTime:   time.Now(),
			Error:     err,
		})
		return err
	}

	ctx.completedAttestors = append(ctx.completedAttestors, CompletedAttestor{
		Attestor:  attestor,
		StartTime: startTime,
		EndTime:   time.Now(),
	})

	if materialer, ok := attestor.(Materialer); ok {
		ctx.addMaterials(materialer)
	}

	if producter, ok := attestor.(Producer); ok {
		ctx.addProducts(producter)
	}

	return nil
}

func (ctx *AttestationContext) CompletedAttestors() []CompletedAttestor {
	attestors := make([]CompletedAttestor, len(ctx.completedAttestors))
	copy(attestors, ctx.completedAttestors)
	return attestors
}

func (ctx *AttestationContext) WorkingDir() string {
	return ctx.workingDir
}

func (ctx *AttestationContext) Hashes() []crypto.Hash {
	hashes := make([]crypto.Hash, len(ctx.hashes))
	copy(hashes, ctx.hashes)
	return hashes
}

func (ctx *AttestationContext) Context() context.Context {
	return ctx.ctx
}

func (ctx *AttestationContext) Materials() map[string]cryptoutil.DigestSet {
	matCopy := make(map[string]cryptoutil.DigestSet)
	for k, v := range ctx.materials {
		matCopy[k] = v
	}

	return matCopy
}

func (ctx *AttestationContext) Products() map[string]Product {
	prodCopy := make(map[string]Product)
	for k, v := range ctx.products {
		prodCopy[k] = v
	}

	return ctx.products
}

func (ctx *AttestationContext) addMaterials(materialer Materialer) {
	newMats := materialer.Materials()
	for k, v := range newMats {
		ctx.materials[k] = v
	}
}

func (ctx *AttestationContext) addProducts(producter Producer) {
	newProds := producter.Products()
	for k, v := range newProds {
		ctx.products[k] = v
	}
}
