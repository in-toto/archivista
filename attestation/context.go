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

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/log"
)

type RunType string

const (
	PreMaterialRunType RunType = "prematerial"
	MaterialRunType    RunType = "material"
	ExecuteRunType     RunType = "execute"
	ProductRunType     RunType = "product"
	PostProductRunType RunType = "postproduct"
)

func runTypeOrder() []RunType {
	return []RunType{PreMaterialRunType, MaterialRunType, ExecuteRunType, ProductRunType, PostProductRunType}
}

func (r RunType) String() string {
	return string(r)
}

type ErrAttestor struct {
	Name    string
	RunType RunType
	Reason  string
}

func (e ErrAttestor) Error() string {
	return fmt.Sprintf("error returned for attestor %s of run type %s: %s", e.Name, e.RunType, e.Reason)
}

type AttestationContextOption func(ctx *AttestationContext)

func WithContext(ctx context.Context) AttestationContextOption {
	return func(actx *AttestationContext) {
		actx.ctx = ctx
	}
}

func WithHashes(hashes []cryptoutil.DigestValue) AttestationContextOption {
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
	hashes             []cryptoutil.DigestValue
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
		hashes:     []cryptoutil.DigestValue{{Hash: crypto.SHA256}, {Hash: crypto.SHA256, GitOID: true}, {Hash: crypto.SHA1, GitOID: true}},
		materials:  make(map[string]cryptoutil.DigestSet),
		products:   make(map[string]Product),
	}

	for _, opt := range opts {
		opt(ctx)
	}

	return ctx, nil
}

func (ctx *AttestationContext) RunAttestors() error {
	attestors := make(map[RunType][]Attestor)
	for _, attestor := range ctx.attestors {
		if attestor.RunType() == "" {
			return ErrAttestor{
				Name:    attestor.Name(),
				RunType: attestor.RunType(),
				Reason:  "attestor run type not set",
			}
		}

		attestors[attestor.RunType()] = append(attestors[attestor.RunType()], attestor)
	}

	order := runTypeOrder()
	for _, k := range order {
		log.Debugf("Starting %s attestors...", k.String())
		for _, att := range attestors[k] {
			log.Infof("Starting %v attestor...", att.Name())
			ctx.runAttestor(att)
		}
	}

	return nil
}

func (ctx *AttestationContext) runAttestor(attestor Attestor) {
	startTime := time.Now()
	if err := attestor.Attest(ctx); err != nil {
		ctx.completedAttestors = append(ctx.completedAttestors, CompletedAttestor{
			Attestor:  attestor,
			StartTime: startTime,
			EndTime:   time.Now(),
			Error:     err,
		})
	}

	ctx.completedAttestors = append(ctx.completedAttestors, CompletedAttestor{
		Attestor:  attestor,
		StartTime: startTime,
		EndTime:   time.Now(),
	})

	if materialer, ok := attestor.(Materialer); ok {
		ctx.addMaterials(materialer)
	}

	if producer, ok := attestor.(Producer); ok {
		ctx.addProducts(producer)
	}
}

func (ctx *AttestationContext) CompletedAttestors() []CompletedAttestor {
	out := make([]CompletedAttestor, len(ctx.completedAttestors))
	copy(out, ctx.completedAttestors)
	return out
}

func (ctx *AttestationContext) WorkingDir() string {
	return ctx.workingDir
}

func (ctx *AttestationContext) Hashes() []cryptoutil.DigestValue {
	hashes := make([]cryptoutil.DigestValue, len(ctx.hashes))
	copy(hashes, ctx.hashes)
	return hashes
}

func (ctx *AttestationContext) Context() context.Context {
	return ctx.ctx
}

func (ctx *AttestationContext) Materials() map[string]cryptoutil.DigestSet {
	out := make(map[string]cryptoutil.DigestSet)
	for k, v := range ctx.materials {
		out[k] = v
	}
	return out
}

func (ctx *AttestationContext) Products() map[string]Product {
	out := make(map[string]Product)
	for k, v := range ctx.products {
		out[k] = v
	}
	return out
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
