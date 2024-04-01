// Copyright 2023 The Witness Contributors
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

package gcp

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
	"testing"

	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/signer/kms"
	"github.com/stretchr/testify/assert"
)

func TestParseReference(t *testing.T) {
	tests := []struct {
		in             string
		wantProjectID  string
		wantLocationID string
		wantKeyRing    string
		wantKeyName    string
		wantKeyVersion string
		wantErr        bool
	}{
		{
			in:             "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			wantProjectID:  "pp",
			wantLocationID: "ll",
			wantKeyRing:    "rr",
			wantKeyName:    "kk",
			wantErr:        false,
		},
		{
			in:             "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk/versions/1",
			wantProjectID:  "pp",
			wantLocationID: "ll",
			wantKeyRing:    "rr",
			wantKeyName:    "kk",
			wantKeyVersion: "1",
			wantErr:        false,
		},
		{
			in:             "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk/cryptoKeyVersions/1",
			wantProjectID:  "pp",
			wantLocationID: "ll",
			wantKeyRing:    "rr",
			wantKeyName:    "kk",
			wantKeyVersion: "1",
			wantErr:        false,
		},
		{
			in:      "gcpkms://projects/p1/p2/locations/l1/l2/keyRings/r1/r2/cryptoKeys/k1",
			wantErr: true,
		},
		{
			in:      "foo://bar",
			wantErr: true,
		},
		{
			in:      "",
			wantErr: true,
		},
		{
			in:      "gcpkms://projects/p1/p2/locations/l1/l2/keyRings/r1/r2/cryptoKeys/k1/versions",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			gotProjectID, gotLocationID, gotKeyRing, gotKeyName, gotKeyVersion, err := parseReference(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProjectID != tt.wantProjectID {
				t.Errorf("parseReference() gotProjectID = %v, want %v", gotProjectID, tt.wantProjectID)
			}
			if gotLocationID != tt.wantLocationID {
				t.Errorf("parseReference() gotLocationID = %v, want %v", gotLocationID, tt.wantLocationID)
			}
			if gotKeyRing != tt.wantKeyRing {
				t.Errorf("parseReference() gotKeyRing = %v, want %v", gotKeyRing, tt.wantKeyRing)
			}
			if gotKeyName != tt.wantKeyName {
				t.Errorf("parseReference() gotKeyName = %v, want %v", gotKeyName, tt.wantKeyName)
			}
			if gotKeyVersion != tt.wantKeyVersion {
				t.Errorf("parseReference() gotKeyVersion = %v, want %v", gotKeyVersion, tt.wantKeyVersion)
			}
		})
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		hash        crypto.Hash
		message     string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "successful sign",
			ref:     "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			hash:    crypto.SHA256,
			message: "foo",
			wantErr: false,
		},
		{
			name:    "SHA-512",
			ref:     "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			hash:    crypto.SHA512,
			message: "foo",
			wantErr: false,
		},
		{
			name:        "bad ref",
			ref:         "blablabla",
			hash:        crypto.SHA256,
			message:     "foo",
			wantErr:     true,
			expectedErr: fmt.Errorf("kms specification should be in the format awskms://[ENDPOINT]/[ID/ALIAS/ARN] (endpoint optional)"),
		},
		{
			name:        "unsupported hash algorithm",
			ref:         "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			hash:        crypto.RIPEMD160,
			message:     "foo",
			wantErr:     true,
			expectedErr: fmt.Errorf(`unsupported hash algorithm: "RIPEMD-160" not in [SHA-256 SHA-384 SHA-512]`),
		},
		{
			name:        "aws ref",
			ref:         "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			hash:        crypto.SHA256,
			message:     "foobarbaz",
			wantErr:     true,
			expectedErr: fmt.Errorf("kms specification should be in the format awskms://[ENDPOINT]/[ID/ALIAS/ARN] (endpoint optional)"),
		},
	}

	for _, tt := range tests {
		fmt.Println("sign test: ", tt.name)
		ctx := context.TODO()
		dig, _, err := cryptoutil.ComputeDigest(bytes.NewReader([]byte(tt.message)), tt.hash, gcpSupportedHashFuncs)
		if tt.wantErr && err != nil {
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		ksp := kms.New(kms.WithRef(tt.ref), kms.WithHash(tt.hash.String()))
		c, err := newFakeGCPClient(context.TODO(), ksp)
		if tt.wantErr && err != nil {
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		s, err := c.sign(ctx, dig, tt.hash, 00)
		if tt.wantErr && err != nil {
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		if s == nil {
			t.Fatal("signature is nil")
		}

		if tt.wantErr {
			t.Fatalf("expected test %s to fail", tt.name)
		}

	}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		hash        crypto.Hash
		mess        []string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "successful verify",
			hash:    crypto.SHA256,
			ref:     "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			mess:    []string{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:    "SHA-512",
			hash:    crypto.SHA512,
			ref:     "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			mess:    []string{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:        "bad ref",
			hash:        crypto.SHA256,
			ref:         "blablabla",
			mess:        []string{"foo", "bar", "baz"},
			wantErr:     true,
			expectedErr: fmt.Errorf("kms specification should be in the format awskms://[ENDPOINT]/[ID/ALIAS/ARN] (endpoint optional)"),
		},
		{
			name:        "unsupported hash algorithm",
			hash:        crypto.RIPEMD160,
			ref:         "gcpkms://projects/pp/locations/ll/keyRings/rr/cryptoKeys/kk",
			mess:        []string{"foo", "bar", "baz"},
			wantErr:     true,
			expectedErr: fmt.Errorf(`unsupported hash algorithm: "RIPEMD-160" not in [SHA-256 SHA-384 SHA-512]`),
		},
		{
			name:        "aws ref",
			hash:        crypto.SHA256,
			ref:         "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			mess:        []string{"foo", "bar", "baz"},
			wantErr:     true,
			expectedErr: fmt.Errorf("kms specification should be in the format awskms://[ENDPOINT]/[ID/ALIAS/ARN] (endpoint optional)"),
		},
	}

	for _, tt := range tests {
		errFound := false
		fmt.Println("verify test: ", tt.name)
		ctx := context.TODO()
		ksp := kms.New(kms.WithRef(tt.ref), kms.WithHash(tt.hash.String()))
		c, err := newFakeGCPClient(context.TODO(), ksp)
		if tt.wantErr && err != nil {
			errFound = true
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		for _, mess := range tt.mess {
			bs, bv, err := createTestKey()
			if err != nil {
				t.Fatal(err)
			}

			dig, _, err := cryptoutil.ComputeDigest(bytes.NewReader([]byte(mess)), tt.hash, gcpSupportedHashFuncs)
			if tt.wantErr && err != nil {
				errFound = true
				assert.ErrorAs(t, err, &tt.expectedErr)
				continue
			} else if err != nil {
				t.Fatal(err)
			}

			sig, err := c.sign(ctx, []byte(dig), crypto.SHA256, 00)
			if tt.wantErr && err != nil {
				errFound = true
				assert.ErrorAs(t, err, &tt.expectedErr)
				continue
			} else if err != nil {
				t.Fatal(err)
			}

			bsig, err := bs.Sign(bytes.NewReader([]byte(dig)))
			if err != nil {
				t.Fatal(err)
			}

			r := bytes.NewReader([]byte(dig))

			err = c.verify(r, sig)
			if tt.wantErr && err != nil {
				errFound = true
				assert.ErrorAs(t, err, &tt.expectedErr)
				continue
			} else if err != nil {
				t.Fatal(err)
			}

			err = c.verify(r, bsig)
			if err == nil {
				t.Fatal("expected verification to fail")
			}

			err = bv.Verify(bytes.NewReader([]byte(dig)), sig)
			if err == nil {
				t.Fatal("expected verification to fail")
			}

		}

		if tt.wantErr && !errFound {
			t.Fatalf("expected test %s to fail with error %s", tt.name, tt.expectedErr.Error())
		}
	}
}
