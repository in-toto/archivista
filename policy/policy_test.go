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

package policy

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"
	"time"

	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/attestation/commandrun"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/intoto"
	"github.com/in-toto/go-witness/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	attestation.RegisterAttestation("dummy-prods", "dummy-prods", attestation.PostProductRunType, func() attestation.Attestor {
		return &DummyProducer{}
	})
	attestation.RegisterAttestation("dummy-mats", "dummy-mats", attestation.PreMaterialRunType, func() attestation.Attestor {
		return &DummyMaterialer{}
	})
}

func createTestKey() (cryptoutil.Signer, cryptoutil.Verifier, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, err
	}

	signer := cryptoutil.NewRSASigner(privKey, crypto.SHA256)
	verifier := cryptoutil.NewRSAVerifier(&privKey.PublicKey, crypto.SHA256)
	keyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: keyBytes})
	if err != nil {
		return nil, nil, nil, err
	}

	return signer, verifier, pemBytes, nil
}

func TestVerify(t *testing.T) {
	_, verifier, pubKeyPem, err := createTestKey()
	require.NoError(t, err)
	keyID, err := verifier.KeyID()
	require.NoError(t, err)
	_, verifier2, pubKeyPem2, err := createTestKey()
	require.NoError(t, err)
	keyID2, err := verifier2.KeyID()
	require.NoError(t, err)
	commandPolicy := []byte(`package test
deny[msg] {
	input.cmd != ["go", "build", "./"]
	msg := "unexpected cmd"
}`)
	exitPolicy := []byte(`package commandrun.exitcode
deny[msg] {
	input.exitcode != 0
	msg := "exitcode not 0"
}`)

	policy := Policy{
		Expires: metav1.NewTime(time.Now().Add(1 * time.Hour)),
		PublicKeys: map[string]PublicKey{
			keyID: {
				KeyID: keyID,
				Key:   pubKeyPem,
			},
			keyID2: {
				KeyID: keyID2,
				Key:   pubKeyPem2,
			},
		},
		Steps: map[string]Step{
			"step1": {
				Name: "step1",
				Functionaries: []Functionary{
					{
						Type:        "PublicKey",
						PublicKeyID: keyID,
					},
				},
				Attestations: []Attestation{
					{
						Type: commandrun.Type,
						RegoPolicies: []RegoPolicy{
							{
								Module: commandPolicy,
								Name:   "expected command",
							},
							{
								Name:   "exited successfully",
								Module: exitPolicy,
							},
						},
					},
				},
			},
		},
	}

	commandRun := commandrun.New()
	commandRun.Cmd = []string{"go", "build", "./"}
	commandRun.ExitCode = 0

	step1Collection := attestation.NewCollection("step1", []attestation.CompletedAttestor{
		{
			Attestor:  commandRun,
			StartTime: time.Now().Add(-1 * time.Minute),
			EndTime:   time.Now(),
			Error:     nil,
		},
	})

	step1CollectionJson, err := json.Marshal(&step1Collection)
	require.NoError(t, err)
	intotoStatement, err := intoto.NewStatement(attestation.CollectionType, step1CollectionJson, map[string]cryptoutil.DigestSet{"dummy": {cryptoutil.DigestValue{Hash: crypto.SHA256}: "dummy"}})
	require.NoError(t, err)

	_, err = policy.Verify(
		context.Background(),
		WithSubjectDigests([]string{"dummy"}),
		WithVerifiedSource(
			newDummyVerifiedSourcer([]source.VerifiedCollection{
				{
					Verifiers: []cryptoutil.Verifier{verifier},
					CollectionEnvelope: source.CollectionEnvelope{
						Statement:  intotoStatement,
						Collection: step1Collection,
						Reference:  "1",
					},
				},
			}),
		),
	)
	assert.NoError(t, err)

	_, err = policy.Verify(
		context.Background(),
		WithSubjectDigests([]string{"dummy"}),
		WithVerifiedSource(
			newDummyVerifiedSourcer([]source.VerifiedCollection{
				{
					Verifiers: []cryptoutil.Verifier{},
					CollectionEnvelope: source.CollectionEnvelope{
						Statement:  intotoStatement,
						Collection: step1Collection,
						Reference:  "1",
					},
				},
			}),
		),
	)
	assert.Error(t, err)
	assert.IsType(t, ErrPolicyDenied{}, err)
}

func TestArtifacts(t *testing.T) {
	_, verifier, pubKeyPem, err := createTestKey()
	require.NoError(t, err)
	keyID, err := verifier.KeyID()
	require.NoError(t, err)

	policy := Policy{
		Expires: metav1.NewTime(time.Now().Add(1 * time.Hour)),
		PublicKeys: map[string]PublicKey{
			keyID: {
				KeyID: keyID,
				Key:   pubKeyPem,
			},
		},
		Steps: map[string]Step{
			"step1": {
				Name: "step1",
				Functionaries: []Functionary{
					{
						Type:        "PublicKey",
						PublicKeyID: keyID,
					},
				},
				Attestations: []Attestation{
					{
						Type: "dummy-prods",
					},
				},
			},
			"step2": {
				Name:          "step2",
				ArtifactsFrom: []string{"step1"},
				Functionaries: []Functionary{
					{
						Type:        "PublicKey",
						PublicKeyID: keyID,
					},
				},
				Attestations: []Attestation{
					{
						Type: "dummy-mats",
					},
				},
			},
		},
	}

	dummySha := "a1073968266a4ed65472a80ebcfd31f1955cfdf8f23d439b1df84d78ce05f7a9"
	path := "testfile"
	mats := map[string]cryptoutil.DigestSet{path: {cryptoutil.DigestValue{Hash: crypto.SHA256}: dummySha}}
	prods := map[string]attestation.Product{path: {Digest: cryptoutil.DigestSet{cryptoutil.DigestValue{Hash: crypto.SHA256}: dummySha}, MimeType: "application/text"}}

	step1Collection := attestation.NewCollection("step1", []attestation.CompletedAttestor{
		{
			Attestor:  DummyProducer{prods},
			StartTime: time.Now().Add(-1 * time.Minute),
			EndTime:   time.Now(),
			Error:     nil,
		},
	})

	step2Collection := attestation.NewCollection("step2", []attestation.CompletedAttestor{
		{
			Attestor:  DummyMaterialer{mats},
			StartTime: time.Now().Add(-1 * time.Minute),
			EndTime:   time.Now(),
			Error:     nil,
		},
	})

	step1CollectionJson, err := json.Marshal(step1Collection)
	require.NoError(t, err)
	step2CollectionJson, err := json.Marshal(step2Collection)
	require.NoError(t, err)
	intotoStatement1, err := intoto.NewStatement(attestation.CollectionType, step1CollectionJson, map[string]cryptoutil.DigestSet{})
	require.NoError(t, err)
	intotoStatement2, err := intoto.NewStatement(attestation.CollectionType, step2CollectionJson, map[string]cryptoutil.DigestSet{})
	require.NoError(t, err)
	_, err = policy.Verify(
		context.Background(),
		WithSubjectDigests([]string{dummySha}),
		WithVerifiedSource(newDummyVerifiedSourcer([]source.VerifiedCollection{
			{
				Verifiers: []cryptoutil.Verifier{verifier},
				CollectionEnvelope: source.CollectionEnvelope{
					Statement:  intotoStatement1,
					Collection: step1Collection,
					Reference:  "1",
				},
			},
			{
				Verifiers: []cryptoutil.Verifier{verifier},
				CollectionEnvelope: source.CollectionEnvelope{
					Statement:  intotoStatement2,
					Collection: step2Collection,
					Reference:  "2",
				},
			},
		})),
	)
	assert.NoError(t, err)

	mats[path][cryptoutil.DigestValue{Hash: crypto.SHA256}] = "badhash"

	step2Collection = attestation.NewCollection("step2", []attestation.CompletedAttestor{
		{
			Attestor:  DummyMaterialer{mats},
			StartTime: time.Now().Add(-1 * time.Minute),
			EndTime:   time.Now(),
			Error:     nil,
		},
	})

	step2CollectionJson, err = json.Marshal(step2Collection)
	require.NoError(t, err)
	intotoStatement2, err = intoto.NewStatement(attestation.CollectionType, step2CollectionJson, map[string]cryptoutil.DigestSet{})
	require.NoError(t, err)
	_, err = policy.Verify(
		context.Background(),
		WithSubjectDigests([]string{dummySha}),
		WithVerifiedSource(newDummyVerifiedSourcer([]source.VerifiedCollection{
			{
				Verifiers: []cryptoutil.Verifier{verifier},
				CollectionEnvelope: source.CollectionEnvelope{
					Statement:  intotoStatement1,
					Collection: step1Collection,
					Reference:  "1",
				},
			},
			{
				Verifiers: []cryptoutil.Verifier{verifier},
				CollectionEnvelope: source.CollectionEnvelope{
					Statement:  intotoStatement2,
					Collection: step2Collection,
					Reference:  "2",
				},
			},
		})),
	)
	assert.Error(t, err)
	assert.IsType(t, ErrPolicyDenied{}, err)
}

type DummyMaterialer struct {
	M map[string]cryptoutil.DigestSet
}

func (DummyMaterialer) Name() string {
	return "dummy-mats"
}

func (DummyMaterialer) Type() string {
	return "dummy-mats"
}

func (DummyMaterialer) RunType() attestation.RunType {
	return attestation.PreMaterialRunType
}

func (DummyMaterialer) Attest(*attestation.AttestationContext) error {
	return nil
}

func (m DummyMaterialer) Materials() map[string]cryptoutil.DigestSet {
	return m.M
}

type DummyProducer struct {
	P map[string]attestation.Product
}

func (DummyProducer) Name() string {
	return "dummy-prods"
}

func (DummyProducer) Type() string {
	return "dummy-prods"
}

func (DummyProducer) RunType() attestation.RunType {
	return attestation.PostProductRunType
}

func (DummyProducer) Attest(*attestation.AttestationContext) error {
	return nil
}

func (m DummyProducer) Products() map[string]attestation.Product {
	return m.P
}

type dummyVerifiedSourcer struct {
	verifiedCollections []source.VerifiedCollection
}

func newDummyVerifiedSourcer(verifiedCollections []source.VerifiedCollection) *dummyVerifiedSourcer {
	return &dummyVerifiedSourcer{verifiedCollections}
}

func (s *dummyVerifiedSourcer) Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]source.VerifiedCollection, error) {
	return s.verifiedCollections, nil
}
