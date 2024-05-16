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

	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/internal/metadatastorage"
	"github.com/openvex/go-vex/pkg/vex"
)

const (
	Predicate = "https://openvex.dev/ns"
)

type ParsedVEX vex.VEX

func Parse(data []byte) (metadatastorage.Storer, error) {
	parsedVex := ParsedVEX{}
	if err := json.Unmarshal(data, &parsedVex); err != nil {
		return parsedVex, err
	}

	return parsedVex, nil
}

func (pv ParsedVEX) Store(ctx context.Context, tx *ent.Tx, stmtID int) error {
	document, err := tx.VexDocument.Create().
		SetStatementID(stmtID).
		SetVexID(pv.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	for _, s := range pv.Statements {
		if err := tx.VexStatement.Create().
			SetVexDocumentID(document.ID).
			SetVulnID(string(s.Vulnerability.Name)).
			Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
