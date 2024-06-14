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

// A note: this follows a pattern followed by network service mesh.
// The pattern was copied from the Network Service Mesh Project
// and modified for use here. The original code was published under the
// Apache License V2.

package omnibor_parser

import (
	"context"
	"encoding/json"
	"github.com/fkautz/omnitrail-go"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage/attestationcollection"
)

func init() {
	// register with parser_registry if the parser_registry exists
	attestationcollection.Register("https://witness.dev/attestations/omnitrail/v0.1", Parse)
}

func Parse(ctx context.Context, tx *ent.Tx, attestation *ent.Attestation, attestationType string, message json.RawMessage) error {
	var envelope struct {
		Envelope json.RawMessage `json:"Envelope"`
	}
	// Unmarshal the attestation into an envelope.
	if err := json.Unmarshal(message, &envelope); err != nil {
		return err
	}

	var omnitrailData omnitrail.Envelope
	// Unmarshal the envelope into omnitrailData.
	if err := json.Unmarshal(envelope.Envelope, &omnitrailData); err != nil {
		return err
	}

	// Create a new Omnitrail entity in the database.
	omnitrailEntity, err := tx.Omnitrail.Create().
		SetAttestationID(attestation.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	// Iterate over each mapping in the omnitrail data.
	for key, element := range omnitrailData.Mapping {
		// Create a new Mapping entity in the database.
		mappingEntity, err := tx.Mapping.Create().
			SetPath(key).
			SetOmnitrailID(omnitrailEntity.ID).
			SetType(element.Type).
			SetSha1(element.Sha1).
			SetSha256(element.Sha256).
			SetGitoidSha1(element.Sha1Gitoid).
			SetGitoidSha256(element.Sha256Gitoid).
			Save(ctx)
		if err != nil {
			return err
		}

		// Extract posix data from the element.
		posixData := element.Posix
		// Create a new Posix entity in the database.
		_, err = tx.Posix.Create().
			SetAtime(posixData.ATime).
			SetCtime(posixData.CTime).
			SetCreationTime(posixData.CreationTime).
			SetExtendedAttributes(posixData.ExtendedAttributes).
			SetFileDeviceID(posixData.FileDeviceID).
			SetFileFlags(posixData.FileFlags).
			SetFileInode(posixData.FileInode).
			SetFileSystemID(posixData.FileSystemID).
			SetFileType(posixData.FileType).
			SetHardLinkCount(posixData.HardLinkCount).
			SetMtime(posixData.MTime).
			SetMetadataCtime(posixData.MetadataCTime).
			SetOwnerUID(posixData.OwnerUID).
			SetOwnerGid(posixData.OwnerGID).
			SetPermissions(posixData.Permissions).
			SetSize(posixData.Size).
			SetMappingID(mappingEntity.ID).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
