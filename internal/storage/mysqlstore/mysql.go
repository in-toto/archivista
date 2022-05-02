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
	"encoding/json"
	"fmt"
	"github.com/git-bom/gitbom-go"
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

func (s *store) GetBySubjectDigest(ctx context.Context, request *archivist.GetBySubjectDigestRequest) (*archivist.GetBySubjectDigestResponse, error) {
	digests, err := s.client.Digest.Query().Where(
		digest.And(
			digest.Algorithm(request.Algorithm),
			digest.Value(request.Value),
		),
	).All(ctx)

	logrus.WithContext(ctx).Printf("digests: %+v", digests)

	results := make([]string, 0)
	for _, curDigest := range digests {
		curDsse, err := curDigest.QuerySubject().QueryStatement().QueryDsse().Only(ctx)
		if err != nil {
			logrus.WithContext(ctx).Errorf("error getting statement: %+v", err)
		}
		results = append(results, curDsse.GitbomSha256)
	}
	logrus.WithContext(ctx).Printf("statements: %+v", results)

	logrus.WithContext(ctx).Printf("results: %+v", results)

	return &archivist.GetBySubjectDigestResponse{Object: results}, err
}

func (s *store) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	tx, err := s.client.Tx(ctx)
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

	// generate gitbom
	gb := gitbom.NewSha256GitBom()
	if err := gb.AddReference([]byte(obj), nil); err != nil {
		logrus.WithContext(ctx).Errorf("gitbom tag generation failed: %+v", err)
		return nil, err
	}

	dsse, err := tx.Dsse.Create().
		SetPayloadType(envelope.PayloadType).
		SetGitbomSha256(gb.Identity()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Statement.Create().
		AddDsse(dsse).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	for _, subject := range payload.Subject {
		storedSubject, err := tx.Subject.Create().
			SetName(subject.Name).
			AddStatement(stmt).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		for algorithm, value := range subject.Digest {
			if err := tx.Digest.Create().
				SetAlgorithm(algorithm).
				SetValue(value).SetSubject(storedSubject).
				Exec(ctx); err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()

	if err != nil {
		logrus.Errorf("unable to commit transaction: %+v", err)
		return nil, err
	}

	fmt.Println("stored!")

	return &emptypb.Empty{}, nil
}
