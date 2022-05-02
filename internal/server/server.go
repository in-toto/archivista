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
func (s *archivistServer) GetBySubjectDigest(ctx context.Context, request *archivist.GetBySubjectDigestRequest) (*archivist.GetBySubjectDigestResponse, error) {
	logrus.WithContext(ctx).Printf("retrieving by subject... ")
	return s.store.GetBySubjectDigest(ctx, request)
}

type collectorServer struct {
	archivist.UnimplementedCollectorServer

	store archivist.CollectorServer
}

func NewCollectorServer(store archivist.CollectorServer) archivist.CollectorServer {
	return &collectorServer{
		store: store,
	}
}

func (s *collectorServer) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	fmt.Println("middleware: store")
	res, err := s.store.Store(ctx, request)
	if err != nil {
		logrus.WithContext(ctx).Printf("received error from database: %+v", err)
		return nil, err
	}
	return res, nil
}
