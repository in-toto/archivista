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
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"google.golang.org/protobuf/types/known/emptypb"
)

type archivistServer struct {
	archivist.UnimplementedArchivistServer

	store archivist.ArchivistServer
}

func NewArchivistServer(store archivist.ArchivistServer) archivist.ArchivistServer {
	return &archivistServer{
		store: store,
	}
}

func (s *archivistServer) GetBySubjectDigest(request *archivist.GetBySubjectDigestRequest, server archivist.Archivist_GetBySubjectDigestServer) error {
	ctx := server.Context()
	logrus.WithContext(ctx).Printf("retrieving by subject... ")
	return s.store.GetBySubjectDigest(request, server)
}

type collectorServer struct {
	archivist.UnimplementedCollectorServer

	metadataStore archivist.CollectorServer
	objectStore   archivist.CollectorServer
}

func NewCollectorServer(metadataStore, objectStore archivist.CollectorServer) archivist.CollectorServer {
	return &collectorServer{
		objectStore:   objectStore,
		metadataStore: metadataStore,
	}
}

func (s *collectorServer) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	fmt.Println("middleware: store")
	if _, err := s.metadataStore.Store(ctx, request); err != nil {
		logrus.WithContext(ctx).Printf("received error from metadata store: %+v", err)
		return nil, err
	}

	if _, err := s.objectStore.Store(ctx, request); err != nil {
		logrus.WithContext(ctx).Printf("received error from object store: %+v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *collectorServer) Get(ctx context.Context, request *archivist.GetRequest) (*archivist.GetResponse, error) {
	return s.objectStore.Get(ctx, request)
}
