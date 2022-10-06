// Copyright 2022 The Archivist Contributors
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

package mysqlstore

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect/sql"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"github.com/testifysec/archivist/ent"
	"github.com/testifysec/go-witness/attestation"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/dsse"
	"github.com/testifysec/go-witness/intoto"

	_ "github.com/go-sql-driver/mysql"
)

// mysql has a limit of 65536 parameters in a single query. each subject has ~2 parameters [statment id and name],
// so we can theoretically jam 65535/2 subjects in a single batch. but we probably want some breathing room just in case.
const subjectBatchSize = 30000

// mysql has a limit of 65536 parameters in a single query. each subject has ~3 parameters [subject id, algo, value],
// so we can theoretically jam 65535/3 subjects in a single batch. but we probably want some breathing room just in case.
const subjectDigestBatchSize = 20000

type Store struct {
	client *ent.Client
}

func New(ctx context.Context, connectionstring string) (*Store, <-chan error, error) {
	drv, err := sql.Open("mysql", connectionstring)
	if err != nil {
		return nil, nil, err
	}
	sqlcommentDrv := sqlcomment.NewDriver(drv,
		sqlcomment.WithDriverVerTag(),
		sqlcomment.WithTags(sqlcomment.Tags{
			sqlcomment.KeyApplication: "archivist",
			sqlcomment.KeyFramework:   "net/http",
		}),
	)

	// TODO make sure these take affect in sqlcommentDrv
	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(3 * time.Minute)

	client := ent.NewClient(ent.Driver(sqlcommentDrv))

	errCh := make(chan error)

	go func() {
		<-ctx.Done()
		err := client.Close()
		if err != nil {
			log.FromContext(ctx).Errorf("error closing database: %+v", err)
		}
		close(errCh)
	}()

	if err := client.Schema.Create(ctx); err != nil {
		log.FromContext(ctx).Fatalf("failed creating schema resources: %v", err)
	}

	return &Store{
		client: client,
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

// attestation.Collection from go-witness will try to parse each of the attestations by calling their factory functions,
// which require the attestations to be registered in the go-witness library.  We don't really care about the actual attestation
// data for the purposes here, so just leave it as a raw message.
type parsedCollection struct {
	attestation.Collection
	Attestations []struct {
		Type        string          `json:"type"`
		Attestation json.RawMessage `json:"attestation"`
	} `json:"attestations"`
}

func (s *Store) Store(ctx context.Context, gitoid string, obj []byte) error {
	envelope := &dsse.Envelope{}
	if err := json.Unmarshal(obj, envelope); err != nil {
		return err
	}

	payloadDigestSet, err := cryptoutil.CalculateDigestSetFromBytes(envelope.Payload, []crypto.Hash{crypto.SHA256})
	if err != nil {
		return err
	}

	payload := &intoto.Statement{}
	if err := json.Unmarshal(envelope.Payload, payload); err != nil {
		return err
	}

	parsedCollection := &parsedCollection{}
	if err := json.Unmarshal(payload.Predicate, parsedCollection); err != nil {
		return err
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
			_, err = tx.Signature.Create().
				SetKeyID(sig.KeyID).
				SetSignature(base64.StdEncoding.EncodeToString(sig.Signature)).
				SetDsse(dsse).
				Save(ctx)
			if err != nil {
				return err
			}
		}

		for hashFn, digest := range payloadDigestSet {
			hashName, err := cryptoutil.HashToString(hashFn)
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

		collection, err := tx.AttestationCollection.Create().
			SetStatementID(stmt.ID).
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
	})

	if err != nil {
		log.FromContext(ctx).Errorf("unable to store metadata: %+v", err)
		return err
	}

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
