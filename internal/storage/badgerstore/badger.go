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

package badgerstore

import (
	"context"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	//"google.golang.org/protobuf/types/known/emptypb"
)

type UnifiedStorage interface {
	archivist.ArchivistServer
	archivist.CollectorServer
}

type store struct {
	archivist.UnimplementedArchivistServer
	archivist.UnimplementedCollectorServer

	db *badger.DB
}

func NewServer(ctx context.Context, file string) (UnifiedStorage, chan error, error) {
	errCh := make(chan error)
	var opt = badger.DefaultOptions("").WithInMemory(true)
	if file == "" {
		opt = badger.DefaultOptions("").WithInMemory(true)
	} else {
		opt = badger.DefaultOptions(file)
	}
	db, err := badger.Open(opt)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		<-ctx.Done()
		err := db.Close()
		if err != nil {
			logrus.WithContext(ctx).Errorf("error closing database: %+v", err)
		}
		close(errCh)
	}()

	return &store{
		db: db,
	}, errCh, nil
}

func (s *store) GetBySubject(_ context.Context, request *archivist.GetBySubjectRequest) (*archivist.GetBySubjectResponse, error) {
	results := make([]string, 0)
	err := s.db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(request.GetSubject())
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				results = append(results, string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &archivist.GetBySubjectResponse{Object: results}, err
}

func (s *store) Store(_ context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	err := s.db.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("subject-%d", rand.Int())
		e := badger.NewEntry([]byte(key), []byte(request.GetObject()))
		err := txn.SetEntry(e)
		return err
	})
	return &emptypb.Empty{}, err
}
