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
	"github.com/git-bom/gitbom-go"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/archivist/ent"
	"github.com/testifysec/archivist/ent/subjectdigest"
	"github.com/testifysec/go-witness/attestation"
	"github.com/testifysec/go-witness/cryptoutil"
	"github.com/testifysec/go-witness/dsse"
	"github.com/testifysec/go-witness/intoto"
	"google.golang.org/protobuf/types/known/emptypb"

	_ "github.com/go-sql-driver/mysql"
)

type UnifiedStorage interface {
	archivist.ArchivistServer
	archivist.CollectorServer
}

type store struct {
	archivist.UnimplementedArchivistServer
	archivist.UnimplementedCollectorServer

	client        *ent.Client
	objectStorage archivist.CollectorServer
}

func NewServer(ctx context.Context, connectionstring string, objectStorage archivist.CollectorServer) (UnifiedStorage, <-chan error, error) {
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
			logrus.WithContext(ctx).Errorf("error closing database: %+v", err)
		}
		close(errCh)
	}()

	if err := client.Schema.Create(ctx); err != nil {
		logrus.WithContext(ctx).Fatalf("failed creating schema resources: %v", err)
	}

	return &store{
		client:        client,
		objectStorage: objectStorage,
	}, errCh, nil
}

func (s *store) GetBySubjectDigest(ctx context.Context, request *archivist.GetBySubjectDigestRequest) (*archivist.GetBySubjectDigestResponse, error) {
	res, err := s.client.SubjectDigest.Query().Where(
		subjectdigest.And(
			subjectdigest.Algorithm(request.Algorithm),
			subjectdigest.Value(request.Value),
		),
	).WithSubject(func(q *ent.SubjectQuery) {
		q.WithStatement(func(q *ent.StatementQuery) {
			q.WithDsse()
		})
	}).All(ctx)

	results := make([]string, 0)
	for _, curDigest := range res {
		for _, curDsse := range curDigest.Edges.Subject.Edges.Statement.Edges.Dsse {
			results = append(results, curDsse.GitbomSha256)
		}
	}

	return &archivist.GetBySubjectDigestResponse{Object: results}, err
}

func (s *store) withTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
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

func (s *store) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	fmt.Println("STORING")

	obj := request.GetObject()
	envelope := &dsse.Envelope{}
	if err := json.Unmarshal([]byte(obj), envelope); err != nil {
		return nil, err
	}

	payloadDigestSet, err := cryptoutil.CalculateDigestSetFromBytes(envelope.Payload, []crypto.Hash{crypto.SHA256})
	if err != nil {
		return nil, err
	}

	payload := &intoto.Statement{}
	if err := json.Unmarshal(envelope.Payload, payload); err != nil {
		return nil, err
	}

	parsedCollection := &parsedCollection{}
	if err := json.Unmarshal(payload.Predicate, parsedCollection); err != nil {
		return nil, err
	}

	// generate gitbom
	gb := gitbom.NewSha256GitBom()
	if err := gb.AddReference([]byte(obj), nil); err != nil {
		logrus.WithContext(ctx).Errorf("gitbom tag generation failed: %+v", err)
		return nil, err
	}

	err = s.withTx(ctx, func(tx *ent.Tx) error {
		dsse, err := tx.Dsse.Create().
			SetPayloadType(envelope.PayloadType).
			SetGitbomSha256(gb.Identity()).
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

		for _, subject := range payload.Subject {
			storedSubject, err := tx.Subject.Create().
				SetName(subject.Name).
				SetStatement(stmt).
				Save(ctx)
			if err != nil {
				return err
			}

			for algorithm, value := range subject.Digest {
				if err := tx.SubjectDigest.Create().
					SetAlgorithm(algorithm).
					SetValue(value).SetSubject(storedSubject).
					Exec(ctx); err != nil {
					return err
				}
			}
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
		logrus.Errorf("unable to store metadata: %+v", err)
		return nil, err
	}

	fmt.Println("metadata stored")

	if s.objectStorage != nil {
		_, err = s.objectStorage.Store(ctx, request)
		if err != nil {
			return nil, err
		}

		fmt.Println("object stored")
	}

	return &emptypb.Empty{}, nil
}
