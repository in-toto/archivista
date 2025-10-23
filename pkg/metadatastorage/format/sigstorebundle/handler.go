// Copyright 2025 The Archivista Contributors
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

package sigstorebundle

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage/format"
	"github.com/in-toto/archivista/pkg/sigstorebundle"
	"github.com/in-toto/go-witness/dsse"
	"github.com/sirupsen/logrus"
)

// Handler processes Sigstore bundle format
type Handler struct{}

func init() {
	format.RegisterHandler(&Handler{})
}

// Detect returns true if obj is a valid Sigstore bundle
func (h *Handler) Detect(obj []byte) bool {
	return sigstorebundle.IsBundleJSON(obj)
}

// Store processes and stores a Sigstore bundle
func (h *Handler) Store(ctx context.Context, store format.Store, gitoid string, obj []byte) error {
	logrus.Infof("detected Sigstore bundle, gitoid: %s", gitoid)

	bundle := &sigstorebundle.Bundle{}
	if err := json.Unmarshal(obj, bundle); err != nil {
		logrus.Warnf("failed to unmarshal bundle: %v", err)
		return err
	}

	logrus.Infof("parsed bundle - mediaType: %s, hasDSSE: %v, hasMsgSig: %v",
		bundle.MediaType, bundle.DsseEnvelope != nil, bundle.MessageSignature != nil)

	// Handle DSSE bundles (convert to DSSE envelope for storage)
	if bundle.DsseEnvelope != nil {
		logrus.Infof("processing DSSE bundle: %s", gitoid)
		envelope, err := sigstorebundle.MapBundleToDSSE(bundle, store.GetBundleLimits())
		if err != nil {
			return fmt.Errorf("failed to convert bundle to DSSE: %w", err)
		}

		// Store bundle metadata along with the DSSE envelope
		return h.storeBundle(ctx, store, bundle, envelope, gitoid)
	}

	// Message signature bundles are not yet supported (would need separate storage)
	if bundle.MessageSignature != nil {
		logrus.Warnf("message signature bundle received (not stored in attestations): %s", gitoid)
		// For now, we'll skip storing these bundles in the attestation metadata store
		// In the future, we can implement support for message signatures
		return nil
	}

	// If we get here, it's a bundle with neither DSSE nor message signature
	return fmt.Errorf("bundle has no content: missing both dsseEnvelope and messageSignature")
}

// storeBundle stores a Sigstore bundle by first storing the DSSE envelope,
// then creating a SigstoreBundle record linking to the DSSE
func (h *Handler) storeBundle(ctx context.Context, store format.Store, bundle *sigstorebundle.Bundle, envelope *dsse.Envelope, gitoid string) error {
	// First, store the DSSE attestation and get its ID
	dsseID, err := store.StoreAttestation(ctx, envelope, gitoid)
	if err != nil {
		return err
	}

	// Then store bundle metadata in a separate transaction
	err = store.WithTx(ctx, func(tx *ent.Tx) error {
		_, err := sigstorebundle.StoreBundleMetadata(ctx, tx, gitoid, bundle.MediaType, dsseID)
		return err
	})

	if err != nil {
		logrus.Errorf("unable to store bundle metadata: %+v", err)
		return err
	}

	logrus.Debugf("Stored Sigstore bundle %s with mediaType %s", gitoid, bundle.MediaType)
	return nil
}
