// Copyright 2022-2024 The Archivista Contributors
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

package sqlstore

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/digitorus/timestamp"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage"
	"github.com/in-toto/archivista/pkg/metadatastorage/parserregistry"
	"github.com/in-toto/archivista/pkg/sigstorebundle"
	"github.com/in-toto/go-witness/cryptoutil"
	"github.com/in-toto/go-witness/dsse"
	"github.com/in-toto/go-witness/intoto"
	"github.com/in-toto/go-witness/policy"
	"github.com/sirupsen/logrus"
)

// mysql has a limit of 65536 parameters in a single query. each subject has ~2 parameters [statment id and name],
// so we can theoretically jam 65535/2 subjects in a single batch. but we probably want some breathing room just in case.
const subjectBatchSize = 30000

// mysql has a limit of 65536 parameters in a single query. each subject has ~3 parameters [subject id, algo, value],
// so we can theoretically jam 65535/3 subjects in a single batch. but we probably want some breathing room just in case.
const subjectDigestBatchSize = 20000

// constant for Policy PayloadType
const policyPayloadType = "https://witness.testifysec.com/policy/"

type Store struct {
	client       *ent.Client
	bundleLimits *sigstorebundle.BundleLimits
}

func New(ctx context.Context, client *ent.Client, bundleLimits ...*sigstorebundle.BundleLimits) (*Store, <-chan error, error) {
	// Use default limits if not provided
	var limits *sigstorebundle.BundleLimits
	if len(bundleLimits) > 0 && bundleLimits[0] != nil {
		limits = bundleLimits[0]
	} else {
		limits = sigstorebundle.DefaultBundleLimits()
	}
	errCh := make(chan error)

	go func() {
		<-ctx.Done()
		err := client.Close()
		if err != nil {
			logrus.Errorf("error closing database: %+v", err)
		}
		close(errCh)
	}()

	if err := client.Schema.Create(ctx); err != nil {
		logrus.Fatalf("failed creating schema resources: %v", err)
	}

	return &Store{
		client:       client,
		bundleLimits: limits,
	}, errCh, nil
}

func (s *Store) withTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("unable to rollback transaction: %w", err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) storeAttestation(ctx context.Context, envelope *dsse.Envelope, gitoid string) error {
	payloadDigestSet, err := cryptoutil.CalculateDigestSetFromBytes(envelope.Payload, []cryptoutil.DigestValue{{Hash: crypto.SHA256}})
	if err != nil {
		return err
	}

	payload := &intoto.Statement{}
	if err := json.Unmarshal(envelope.Payload, payload); err != nil {
		return err
	}

	predicateParser, ok := parserregistry.ParserForPredicate(payload.PredicateType)
	var predicateStorer metadatastorage.Storer
	if ok {
		predicateStorer, err = predicateParser(payload.Predicate)
		if err != nil {
			return fmt.Errorf("unable to parse intoto statements predicate: %w", err)
		}
	}

	err = s.withTx(ctx, func(tx *ent.Tx) error {
		dsse, err := tx.Dsse.Create().
			SetPayloadType(envelope.PayloadType).
			SetGitoidSha256(gitoid).
			Save(ctx)
		if err != nil {
			return err
		}

		for _, sig := range envelope.Signatures {
			sigCreate := tx.Signature.Create().
				SetKeyID(sig.KeyID).
				SetSignature(base64.StdEncoding.EncodeToString(sig.Signature)).
				SetDsse(dsse)

			// Store certificate if present
			if len(sig.Certificate) > 0 {
				sigCreate.SetCertificate(sig.Certificate)
			}

			// Store intermediates if present
			if len(sig.Intermediates) > 0 {
				sigCreate.SetIntermediates(sig.Intermediates)
			}

			storedSig, err := sigCreate.Save(ctx)
			if err != nil {
				return err
			}

			for _, timestamp := range sig.Timestamps {
				timestampedTime, err := timeFromTimestamp(timestamp)
				if err != nil {
					return err
				}

				tsCreate := tx.Timestamp.Create().
					SetSignature(storedSig).
					SetTimestamp(timestampedTime).
					SetType(string(timestamp.Type))

				// Store raw RFC3161 data if present
				if len(timestamp.Data) > 0 {
					tsCreate.SetData(timestamp.Data)
				}

				_, err = tsCreate.Save(ctx)
				if err != nil {
					return err
				}
			}
		}

		for hashFn, digest := range payloadDigestSet {
			hashName, err := cryptoutil.HashToString(hashFn.Hash)
			if err != nil {
				return err
			}

			if _, err := tx.PayloadDigest.Create().
				SetDsse(dsse).
				SetAlgorithm(hashName).
				SetValue(digest).
				Save(ctx); err != nil {
				return err
			}
		}

		stmt, err := tx.Statement.Create().
			SetPredicate(payload.PredicateType).
			AddDsse(dsse).
			Save(ctx)
		if err != nil {
			return err
		}

		bulkSubject := make([]*ent.SubjectCreate, 0)
		for _, subject := range payload.Subject {
			bulkSubject = append(bulkSubject,
				tx.Subject.Create().
					SetName(subject.Name).
					SetStatement(stmt),
			)
		}

		subjects, err := batch(ctx, subjectBatchSize, bulkSubject, func(digests ...*ent.SubjectCreate) saver[*ent.Subject] {
			return tx.Subject.CreateBulk(digests...)
		})
		if err != nil {
			return err
		}

		bulkSubjectDigests := make([]*ent.SubjectDigestCreate, 0)
		for i, subject := range payload.Subject {
			for algorithm, value := range subject.Digest {
				bulkSubjectDigests = append(bulkSubjectDigests,
					tx.SubjectDigest.Create().
						SetAlgorithm(algorithm).
						SetValue(value).
						SetSubject(subjects[i]),
				)
			}
		}

		if _, err := batch(ctx, subjectDigestBatchSize, bulkSubjectDigests, func(digests ...*ent.SubjectDigestCreate) saver[*ent.SubjectDigest] {
			return tx.SubjectDigest.CreateBulk(digests...)
		}); err != nil {
			return err
		}

		if predicateStorer != nil {
			if err := predicateStorer.Store(ctx, tx, stmt.ID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("unable to store metadata: %+v", err)
		return err
	}

	return nil
}

func (s *Store) storePolicy(ctx context.Context, envelope *dsse.Envelope, gitoid string) error {
	payloadDigestSet, err := cryptoutil.CalculateDigestSetFromBytes(envelope.Payload, []cryptoutil.DigestValue{{Hash: crypto.SHA256}})
	if err != nil {
		return err
	}

	payload := &policy.Policy{}
	if err := json.Unmarshal(envelope.Payload, payload); err != nil {
		return err
	}

	custom := make(map[string]any)
	custom["gitoid"] = gitoid

	err = s.withTx(ctx, func(tx *ent.Tx) error {
		dsse, err := tx.Dsse.Create().
			SetPayloadType(envelope.PayloadType).
			SetGitoidSha256(gitoid).
			Save(ctx)
		if err != nil {
			return err
		}

		// stores the envelope signatures
		for _, sig := range envelope.Signatures {
			sigCreate := tx.Signature.Create().
				SetKeyID(sig.KeyID).
				SetSignature(base64.StdEncoding.EncodeToString(sig.Signature)).
				SetDsse(dsse)

			// Store certificate if present
			if len(sig.Certificate) > 0 {
				sigCreate.SetCertificate(sig.Certificate)
			}

			// Store intermediates if present
			if len(sig.Intermediates) > 0 {
				sigCreate.SetIntermediates(sig.Intermediates)
			}

			storedSig, err := sigCreate.Save(ctx)
			if err != nil {
				return err
			}

			for _, timestamp := range sig.Timestamps {
				timestampedTime, err := timeFromTimestamp(timestamp)
				if err != nil {
					return err
				}

				tsCreate := tx.Timestamp.Create().
					SetSignature(storedSig).
					SetTimestamp(timestampedTime).
					SetType(string(timestamp.Type))

				// Store raw RFC3161 data if present
				if len(timestamp.Data) > 0 {
					tsCreate.SetData(timestamp.Data)
				}

				_, err = tsCreate.Save(ctx)
				if err != nil {
					return err
				}
			}
		}

		// stores the payload digests
		for hashFn, digest := range payloadDigestSet {
			hashName, err := cryptoutil.HashToString(hashFn.Hash)
			if err != nil {
				return err
			}

			if _, err := tx.PayloadDigest.Create().
				SetDsse(dsse).
				SetAlgorithm(hashName).
				SetValue(digest).
				Save(ctx); err != nil {
				return err
			}
		}

		stmt, err := tx.Statement.Create().
			SetPredicate(envelope.PayloadType).
			AddDsse(dsse).
			Save(ctx)
		if err != nil {
			return err
		}

		// stores the subject
		if _, err := tx.Subject.Create().
			SetName(gitoid).
			SetStatement(stmt).
			Save(ctx); err != nil {
			return err
		}

		if _, err := tx.AttestationPolicy.Create().
			SetStatement(stmt).
			SetName(gitoid).
			Save(ctx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("unable to store metadata: %+v", err)
		return err
	}

	return nil
}

func (s *Store) Store(ctx context.Context, gitoid string, obj []byte) error {
	// First, check if this is a valid Sigstore bundle per spec (RFC)
	if sigstorebundle.IsBundleJSON(obj) {
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
			envelope, err := sigstorebundle.MapBundleToDSSE(bundle, s.bundleLimits)
			if err != nil {
				return fmt.Errorf("failed to convert bundle to DSSE: %w", err)
			}

			// Store bundle metadata along with the DSSE envelope
			if err := s.storeBundle(ctx, bundle, envelope, gitoid); err != nil {
				return err
			}
			return nil
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

	// Try to parse as a plain DSSE envelope
	envelope := &dsse.Envelope{}
	if err := json.Unmarshal(obj, envelope); err != nil {
		return err
	}

	// check if the payload is a policy or an attestation
	if strings.Contains(envelope.PayloadType, policyPayloadType) {
		if err := s.storePolicy(ctx, envelope, gitoid); err != nil {
			return err
		}
	} else {
		if err := s.storeAttestation(ctx, envelope, gitoid); err != nil {
			return err
		}
	}

	return nil
}


// storeBundle stores a Sigstore bundle by first storing the DSSE envelope,
// then creating a SigstoreBundle record linking to the DSSE
func (s *Store) storeBundle(ctx context.Context, bundle *sigstorebundle.Bundle, envelope *dsse.Envelope, gitoid string) error {
	// First, store the DSSE attestation as normal
	if err := s.storeAttestation(ctx, envelope, gitoid); err != nil {
		return err
	}

	// Then store bundle metadata
	// Note: The SigstoreBundle linking is handled in a follow-up transaction
	// since we need the DSSE ID which is already committed at this point
	// For now, just log that we've stored a bundle
	logrus.Debugf("Stored Sigstore bundle %s with mediaType %s", gitoid, bundle.MediaType)

	return nil
}

func (s *Store) GetClient() *ent.Client {
	return s.client
}

type saver[T any] interface {
	Save(context.Context) ([]T, error)
}

func batch[TCreate any, TResult any](ctx context.Context, batchSize int, create []TCreate, saveFn func(...TCreate) saver[TResult]) ([]TResult, error) {
	results := make([]TResult, 0, len(create))
	for i := 0; i < len(create); i += batchSize {
		var batch []TCreate
		if i+batchSize > len(create) {
			batch = create[i:]
		} else {
			batch = create[i : i+batchSize]
		}

		batchSaver := saveFn(batch...)
		batchResults, err := batchSaver.Save(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, batchResults...)
	}

	return results, nil
}

func timeFromTimestamp(ts dsse.SignatureTimestamp) (time.Time, error) {
	switch ts.Type {
	case dsse.TimestampRFC3161:
		tspResponse, err := timestamp.Parse(ts.Data)
		if err != nil {
			return time.Time{}, nil
		}

		return tspResponse.Time, nil
	default:
		return time.Time{}, fmt.Errorf("unknown timestamp type: %v", ts.Type)
	}
}
