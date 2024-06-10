// Copyright 2024 The Archivista Contributors
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

package sarif_parser

import (
	"context"
	"encoding/json"

	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage/attestationcollection"
	"github.com/in-toto/go-witness/attestation/sarif"
)

func init() {
	attestationcollection.Register("https://witness.dev/attestations/sarif/v0.1", Parse)
}

type ParsedSarif sarif.Attestor

func Parse(ctx context.Context, tx *ent.Tx, attestation *ent.Attestation, attestationType string, message json.RawMessage) error {
	sarifAttestation := sarif.Attestor{}
	if err := json.Unmarshal(message, &sarifAttestation); err != nil {
		return err
	}

	stmt, err := tx.Statement.Create().
		SetPredicate(attestationType).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Sarif.Create().
		SetStatement(stmt).
		SetReportFileName(sarifAttestation.ReportFile).
		Save(ctx)
	if err != nil {
		return err
	}

	for _, r := range sarifAttestation.Report.Runs {
		for _, ru := range r.Tool.Driver.Rules {
			if err := tx.SarifRule.Create().
				SetRuleName(*ru.Name).
				SetRuleID(ru.ID).
				SetShortDescription(*ru.ShortDescription.Text).
				Exec(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}
