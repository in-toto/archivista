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

package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/archivist/internal/storage/blob"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type archivistServer struct {
	archivist.UnimplementedArchivistServer
	indexer blob.Indexer
	store   archivist.ArchivistServer
}

func NewArchivistServer(store archivist.ArchivistServer, indexer blob.Indexer) archivist.ArchivistServer {
	return &archivistServer{
		store:   store,
		indexer: indexer,
	}
}
func (s *archivistServer) GetBySubjectDigest(ctx context.Context, request *archivist.GetBySubjectDigestRequest) (*archivist.GetBySubjectDigestResponse, error) {
	logrus.WithContext(ctx).Printf("retrieving by subject... ")
	resp, err := s.store.GetBySubjectDigest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subject by digest: %v", err)
	}

	shas := resp.GetObject()
	logrus.WithContext(ctx).Printf("shas fetched: %s", shas)
	attestations := make([]string, 0)

	for _, gitbomSha := range shas {
		obj, err := s.indexer.GetBlob(gitbomSha)
		if err != nil {
			return nil, fmt.Errorf("failed fetching ref by %s from store: %v", gitbomSha, err)
		}
		attestations = append(attestations, strings.TrimSpace(string(bytes.Trim(obj, "\x00"))))
	}
	return &archivist.GetBySubjectDigestResponse{Object: attestations}, nil
}

type collectorServer struct {
	archivist.UnimplementedCollectorServer
	indexer blob.Indexer
	store   archivist.CollectorServer
}

func NewCollectorServer(store archivist.CollectorServer, indexer blob.Indexer) archivist.CollectorServer {
	return &collectorServer{
		store:   store,
		indexer: indexer,
	}
}

// Store stores the dsse envelope and its relationships into the backend stores
func (s *collectorServer) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	fmt.Println("middleware: store")

	res, err := s.store.Store(ctx, request)
	if err != nil {
		logrus.WithContext(ctx).Printf("received error from database: %+v", err)
		return nil, err
	}

	envBytes := []byte(request.Object)
	ref, err := s.indexer.GetRef(envBytes)
	if err != nil {
		logrus.WithContext(ctx).Printf("failed to get ref for envelope: %v", err)
		return nil, err
	}

	err = s.indexer.PutBlob(ref, envBytes)
	if err != nil {
		logrus.WithContext(ctx).Printf("failed to put blob in store for envelope: %v", err)
		return nil, err
	}

	return res, nil
}
