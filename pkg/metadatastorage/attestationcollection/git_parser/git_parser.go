// Copyright 2024 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package git_parser

import (
	"context"
	"encoding/json"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage/attestationcollection"
	"github.com/in-toto/go-witness/attestation/git"
)

func init() {
	// register with parser_registry if the parser_registry exists
	attestationcollection.Register("https://witness.dev/attestations/git/v0.1", Parse)
}

// register with parser_registry if the parser_registry exists

func Parse(ctx context.Context, tx *ent.Tx, attestation *ent.Attestation, attestationType string, message json.RawMessage) error {
	gitAttestation := git.Attestor{
		CommitHash:     "",
		Author:         "",
		AuthorEmail:    "",
		CommitterName:  "",
		CommitterEmail: "",
		CommitDate:     "",
		CommitMessage:  "",
		Status:         nil,
		CommitDigest:   nil,
		Signature:      "",
		ParentHashes:   nil,
		TreeHash:       "",
		Refs:           nil,
		Remotes:        nil,
		Tags:           nil,
	}

	if err := json.Unmarshal(message, &gitAttestation); err != nil {
		return err
	}

	if _, err := tx.GitAttestation.Create().
		SetAttestation(attestation).
		SetCommitHash(gitAttestation.CommitHash).
		SetAuthor(gitAttestation.Author).
		SetAuthorEmail(gitAttestation.AuthorEmail).
		SetCommitterName(gitAttestation.CommitterName).
		SetCommitterEmail(gitAttestation.CommitterEmail).
		SetCommitMessage(gitAttestation.CommitMessage).
		SetCommitDate(gitAttestation.CommitDate).
		SetSignature(gitAttestation.Signature).
		SetParentHashes(gitAttestation.ParentHashes).
		SetTreeHash(gitAttestation.TreeHash).
		SetRefs(gitAttestation.Refs).
		SetRemotes(gitAttestation.Remotes).
		//SetTags(gitAttestation.Tags). // Implement after we have a graphql marshal/unmarshal for tags
		Save(ctx); err != nil {
		return err
	}

	return nil
}
