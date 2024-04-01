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

package aws

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
		in           string
		wantEndpoint string
		wantKeyID    string
		wantAlias    string
		wantErr      bool
	}{
		{
			in:           "awskms:///1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "1234abcd-12ab-34cd-56ef-1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			// multi-region key
			in:           "awskms:///mrk-1234abcd12ab34cd56ef1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "mrk-1234abcd12ab34cd56ef1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms:///1234ABCD-12AB-34CD-56EF-1234567890AB",
			wantEndpoint: "",
			wantKeyID:    "1234ABCD-12AB-34CD-56EF-1234567890AB",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms://localhost:4566/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "localhost:4566",
			wantKeyID:    "1234abcd-12ab-34cd-56ef-1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms:///arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms://localhost:4566/arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "localhost:4566",
			wantKeyID:    "arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms:///alias/ExampleAlias",
			wantEndpoint: "",
			wantKeyID:    "alias/ExampleAlias",
			wantAlias:    "alias/ExampleAlias",
			wantErr:      false,
		},
		{
			in:           "awskms://localhost:4566/alias/ExampleAlias",
			wantEndpoint: "localhost:4566",
			wantKeyID:    "alias/ExampleAlias",
			wantAlias:    "alias/ExampleAlias",
			wantErr:      false,
		},
		{
			in:           "awskms:///arn:aws:kms:us-east-2:111122223333:alias/ExampleAlias",
			wantEndpoint: "",
			wantKeyID:    "arn:aws:kms:us-east-2:111122223333:alias/ExampleAlias",
			wantAlias:    "alias/ExampleAlias",
			wantErr:      false,
		},
		{
			in:           "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:alias/ExampleAlias",
			wantEndpoint: "",
			wantKeyID:    "arn:aws-us-gov:kms:us-gov-west-1:111122223333:alias/ExampleAlias",
			wantAlias:    "alias/ExampleAlias",
			wantErr:      false,
		},
		{
			in:           "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			wantAlias:    "",
			wantErr:      false,
		},
		{
			in:           "awskms://localhost:4566/arn:aws:kms:us-east-2:111122223333:alias/ExampleAlias",
			wantEndpoint: "localhost:4566",
			wantKeyID:    "arn:aws:kms:us-east-2:111122223333:alias/ExampleAlias",
			wantAlias:    "alias/ExampleAlias",
			wantErr:      false,
		},
		{
			// missing alias/ prefix
			in:           "awskms:///missingalias",
			wantEndpoint: "",
			wantKeyID:    "",
			wantAlias:    "",
			wantErr:      true,
		},
		{
			// invalid UUID
			in:           "awskms:///1234abcd-12ab-YYYY-56ef-1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "",
			wantAlias:    "",
			wantErr:      true,
		},
		{
			// Currently, references without endpoints must use 3
			// slashes. It would be nice to support this format,
			// but that would be harder to parse.
			in:           "awskms://1234abcd-12ab-34cd-56ef-1234567890ab",
			wantEndpoint: "",
			wantKeyID:    "",
			wantAlias:    "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			gotEndpoint, gotKeyID, gotAlias, err := ParseReference(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEndpoint != tt.wantEndpoint {
				t.Errorf("ParseReference() gotEndpoint = %v, want %v", gotEndpoint, tt.wantEndpoint)
			}
			if gotKeyID != tt.wantKeyID {
				t.Errorf("ParseReference() gotKeyID = %v, want %v", gotKeyID, tt.wantKeyID)
			}
			if gotAlias != tt.wantAlias {
				t.Errorf("ParseReference() gotAlias = %v, want %v", gotAlias, tt.wantAlias)
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
			ref:     "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			hash:    crypto.SHA256,
			message: "foo",
			wantErr: false,
		},
		{
			name:    "SHA-512",
			ref:     "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
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
			ref:         "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			hash:        crypto.RIPEMD160,
			message:     "foo",
			wantErr:     true,
			expectedErr: fmt.Errorf(`unsupported hash algorithm: "RIPEMD-160" not in [SHA-256 SHA-384 SHA-512]`),
		},
		{
			name:        "gcp ref",
			ref:         "gcpkms://projects/testproject-23231/locations/europe-west2/keyRings/test/cryptoKeys/test",
			hash:        crypto.SHA256,
			message:     "foobarbaz",
			wantErr:     true,
			expectedErr: fmt.Errorf("kms specification should be in the format awskms://[ENDPOINT]/[ID/ALIAS/ARN] (endpoint optional)"),
		},
	}

	for _, tt := range tests {
		fmt.Println("sign test: ", tt.name)
		ctx := context.TODO()
		dig, _, err := cryptoutil.ComputeDigest(bytes.NewReader([]byte(tt.message)), tt.hash, awsSupportedHashFuncs)
		if tt.wantErr && err != nil {
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		ksp := kms.New(kms.WithRef(tt.ref), kms.WithHash(tt.hash.String()))
		c, err := newFakeAWSClient(context.TODO(), ksp)
		if tt.wantErr && err != nil {
			assert.ErrorAs(t, err, &tt.expectedErr)
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		s, err := c.sign(ctx, dig, tt.hash)
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
			name:    "successful sign",
			hash:    crypto.SHA256,
			ref:     "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			mess:    []string{"foo", "bar", "baz"},
			wantErr: false,
		},
		{
			name:    "SHA-512",
			hash:    crypto.SHA512,
			ref:     "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
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
			ref:         "awskms:///arn:aws-us-gov:kms:us-gov-west-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			mess:        []string{"foo", "bar", "baz"},
			wantErr:     true,
			expectedErr: fmt.Errorf(`unsupported hash algorithm: "RIPEMD-160" not in [SHA-256 SHA-384 SHA-512]`),
		},
		{
			name:        "gcp ref",
			hash:        crypto.SHA256,
			ref:         "gcpkms://projects/testproject-23231/locations/europe-west2/keyRings/test/cryptoKeys/test",
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
		c, err := newFakeAWSClient(context.TODO(), ksp)
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

			dig, _, err := cryptoutil.ComputeDigest(bytes.NewReader([]byte(mess)), tt.hash, awsSupportedHashFuncs)
			if tt.wantErr && err != nil {
				errFound = true
				assert.ErrorAs(t, err, &tt.expectedErr)
				continue
			} else if err != nil {
				t.Fatal(err)
			}

			sig, err := c.sign(ctx, []byte(dig), crypto.SHA256)
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

			err = c.verifyRemotely(ctx, sig, dig)
			if tt.wantErr && err != nil {
				errFound = true
				assert.ErrorAs(t, err, &tt.expectedErr)
				continue
			} else if err != nil {
				t.Fatal(err)
			}

			err = c.verifyRemotely(ctx, bsig, dig)
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
