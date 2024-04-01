// Copyright 2024 The Witness Contributors
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

package witness

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/in-toto/go-witness/internal/test"
	"github.com/in-toto/go-witness/intoto"
	"github.com/in-toto/go-witness/timestamp"
)

func TestVerifyPolicySignature(t *testing.T) {
	// we dont care about the content of th envelope for this test
	rsaSigner, rsaVerifier, _, err := test.CreateTestKey()
	if err != nil {
		t.Fatal(err)
	}

	badRootCert, _, err := test.CreateRoot()
	if err != nil {
		t.Fatal(err)
	}

	rootCert, key, err := test.CreateRoot()
	if err != nil {
		t.Fatal(err)
	}

	leafCert, leafPriv, err := test.CreateLeaf(rootCert, key)
	if err != nil {
		t.Fatal(err)
	}

	x509Signer, err := cryptoutil.NewSigner(leafPriv, cryptoutil.SignWithCertificate(leafCert))
	if err != nil {
		t.Fatal(err)
	}

	timestampers := []timestamp.FakeTimestamper{
		{T: time.Now()},
		{T: time.Now().Add(12 * time.Hour)},
	}

	// Define the test cases.
	tests := []struct {
		name            string
		signer          cryptoutil.Signer
		verifier        cryptoutil.Verifier
		timestampers    []timestamp.FakeTimestamper
		Roots           []*x509.Certificate
		Intermediates   []*x509.Certificate
		certConstraints VerifyOption
		wantErr         bool
	}{
		{
			name:     "valid rsa signature",
			signer:   rsaSigner,
			verifier: rsaVerifier,
			// passing in timestampers to ensure that it is ignored
			timestampers: timestampers,
			wantErr:      false,
		},
		{
			name:    "invalid rsa signature",
			signer:  rsaSigner,
			Roots:   []*x509.Certificate{rootCert},
			wantErr: true,
		},
		{
			name:   "valid x509 signature",
			signer: x509Signer,
			// We're going to pass in to ensure that it is ignored
			Roots:   []*x509.Certificate{rootCert},
			wantErr: false,
		},
		{
			name:   "valid x509 signature w/ constraints",
			signer: x509Signer,
			// We're going to pass in to ensure that it is ignored
			Roots:           []*x509.Certificate{rootCert},
			certConstraints: VerifyWithPolicyCertConstraints(leafCert.Subject.CommonName, leafCert.DNSNames, []string{"*"}, []string{"*"}, []string{"*"}),
			timestampers:    timestampers,
			wantErr:         false,
		},
		{
			name:   "valid x509 signature w/ bad constraints",
			signer: x509Signer,
			// We're going to pass in to ensure that it is ignored
			Roots:           []*x509.Certificate{rootCert},
			certConstraints: VerifyWithPolicyCertConstraints("foo", []string{"bar"}, []string{"baz"}, []string{"qux"}, []string{"quux"}),
			wantErr:         true,
		},
		{
			name:   "unknown root",
			signer: x509Signer,
			// We're going to pass in to ensure that it is ignored
			Roots:   []*x509.Certificate{badRootCert},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		var ts []timestamp.Timestamper
		for _, t := range tt.timestampers {
			ts = append(ts, t)
		}

		env, err := dsse.Sign(intoto.PayloadType, bytes.NewReader([]byte("this is some test data")), dsse.SignWithTimestampers(ts...), dsse.SignWithSigners(tt.signer))
		if err != nil {
			t.Fatal(err)
		}

		var tv []timestamp.TimestampVerifier
		for _, t := range tt.timestampers {
			tv = append(tv, t)
		}

		vo := verifyOptions{
			policyEnvelope:             env,
			policyVerifiers:            []cryptoutil.Verifier{tt.verifier},
			policyCARoots:              tt.Roots,
			policyTimestampAuthorities: tv,
			policyCommonName:           "*",
			policyDNSNames:             []string{"*"},
			policyOrganizations:        []string{"*"},
			policyURIs:                 []string{"*"},
			policyEmails:               []string{"*"},
		}

		if tt.certConstraints != nil {
			tt.certConstraints(&vo)
		}

		err = verifyPolicySignature(context.TODO(), vo)
		if err != nil && !tt.wantErr {
			t.Errorf("testName = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr)
		} else {
			fmt.Printf("test %s passed\n", tt.name)
		}

	}
}
