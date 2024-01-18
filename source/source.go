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

package source

import (
	"context"
	"encoding/json"

	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/dsse"
	"github.com/in-toto/go-witness/intoto"
)

type CollectionEnvelope struct {
	Envelope   dsse.Envelope
	Statement  intoto.Statement
	Collection attestation.Collection
	Reference  string
}

type Sourcer interface {
	Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]CollectionEnvelope, error)
}

func envelopeToCollectionEnvelope(reference string, env dsse.Envelope) (CollectionEnvelope, error) {
	statement := intoto.Statement{}
	if err := json.Unmarshal(env.Payload, &statement); err != nil {
		return CollectionEnvelope{}, err
	}

	collection := attestation.Collection{}
	if err := json.Unmarshal(statement.Predicate, &collection); err != nil {
		return CollectionEnvelope{}, err
	}

	return CollectionEnvelope{
		Reference:  reference,
		Envelope:   env,
		Statement:  statement,
		Collection: collection,
	}, nil
}
