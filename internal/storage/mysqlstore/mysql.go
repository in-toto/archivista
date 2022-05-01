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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/archivist/ent"
	"github.com/testifysec/archivist/ent/digest"
	"github.com/testifysec/witness/pkg/dsse"
	"github.com/testifysec/witness/pkg/intoto"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	//"google.golang.org/protobuf/types/known/emptypb"

	"entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql"
)

type UnifiedStorage interface {
	archivist.ArchivistServer
	archivist.CollectorServer
}

type store struct {
	archivist.UnimplementedArchivistServer
	archivist.UnimplementedCollectorServer

	client *ent.Client
}

func NewServer(ctx context.Context, connectionstring string) (UnifiedStorage, chan error, error) {
	//time.Sleep(1 * time.Hour)
	//b, _ := ioutil.ReadFile("/etc/hosts")
	//logrus.Fatalln(string(b))
	drv, err := sql.Open("mysql", "root:example@tcp(db)/testify")
	if err != nil {
		return nil, nil, err
	}

	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(3 * time.Minute)

	client := ent.NewClient(ent.Driver(drv))

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
		client: client,
	}, errCh, nil
}

func (s *store) GetBySubject(ctx context.Context, request *archivist.GetBySubjectRequest) (*archivist.GetBySubjectResponse, error) {
	digests, err := s.client.Digest.Query().Where(
		digest.And(
			digest.Algorithm("algo"),
			digest.Value("value"),
		),
	).All(ctx)

	statements := make([]*ent.Statement, 0)
	for _, curDigest := range digests {
		curStatement, err := curDigest.QuerySubject().QueryStatement().Only(ctx)
		if err != nil {
			logrus.WithContext(ctx).Errorf("error getting statement: %+v", err)
		}
		statements = append(statements, curStatement)
	}

	results := make([]string, 0)
	for _, stmt := range statements {
		results = append(results, stmt.Statement)
	}

	return &archivist.GetBySubjectResponse{Object: results}, err
}

func (s *store) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	fmt.Println("STORING")
	obj := request.GetObject()

	envelope := &dsse.Envelope{}

	if err := json.Unmarshal([]byte(obj), envelope); err != nil {
		return nil, err
	}

	payload := &intoto.Statement{}

	if err := json.Unmarshal(envelope.Payload, payload); err != nil {
		return nil, err
	}

	payloadHash := sha256.Sum256(envelope.Payload)
	payloasHashString := base64.URLEncoding.EncodeToString(payloadHash[:])

	stmt, err := s.client.Statement.Create().
		SetStatement(payloasHashString).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	for _, subject := range payload.Subject {
		storedSubject, err := s.client.Subject.Create().
			SetName(subject.Name).
			AddStatement(stmt).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		for algorithm, value := range subject.Digest {
			if err := s.client.Digest.Create().
				SetAlgorithm(algorithm).
				SetValue(value).SetSubject(storedSubject).
				Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	fmt.Println("stored!")

	return &emptypb.Empty{}, nil
}
