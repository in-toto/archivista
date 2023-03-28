// Copyright 2023 The Archivista Contributors
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

package attestationcollection

import (
	"context"
	"encoding/json"

	"github.com/testifysec/archivista/ent"
	"github.com/testifysec/archivista/internal/metadatastorage"
	"github.com/testifysec/go-witness/attestation"
)

const (
	Predicate = "https://witness.testifysec.com/attestation-collection/v0.1"
)

// attestation.Collection from go-witness will try to parse each of the attestations by calling their factory functions,
// which require the attestations to be registered in the go-witness library.  We don't really care about the actual attestation
// data for the purposes here, so just leave it as a raw message.
type ParsedCollection struct {
	attestation.Collection
	Attestations []struct {
		Type        string          `json:"type"`
		Attestation json.RawMessage `json:"attestation"`
	} `json:"attestations"`
}

func Parse(data []byte) (metadatastorage.Storer, error) {
	parsedCollection := ParsedCollection{}
	if err := json.Unmarshal(data, &parsedCollection); err != nil {
		return parsedCollection, err
	}

	return parsedCollection, nil
}

func (parsedCollection ParsedCollection) Store(ctx context.Context, tx *ent.Tx, stmtID int) error {
	collection, err := tx.AttestationCollection.Create().
		SetStatementID(stmtID).
		SetName(parsedCollection.Name).
		Save(ctx)
	if err != nil {
		return err
	}

	for _, a := range parsedCollection.Attestations {
		if err := tx.Attestation.Create().
			SetAttestationCollectionID(collection.ID).
			SetType(a.Type).
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
