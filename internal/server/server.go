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
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/edwarnicke/gitoid"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
)

const ChunkSize = 64 * 1024 //64kb seems to be somewhat of an agreed upon message size when streaming: https://github.com/grpc/grpc.github.io/issues/371

type archivistServer struct {
	archivist.UnimplementedArchivistServer

	store MetadataStorer
}

func NewArchivistServer(store MetadataStorer) archivist.ArchivistServer {
	return &archivistServer{
		store: store,
	}
}

func (s *archivistServer) GetBySubjectDigest(request *archivist.GetBySubjectDigestRequest, server archivist.Archivist_GetBySubjectDigestServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	log.FromContext(ctx).Info("retrieving by subject... ")
	responses, err := s.store.GetBySubjectDigest(ctx, request)
	if err != nil {
		return err
	}

	for response := range responses {
		if err := server.Send(response); err != nil {
			return err
		}
	}

	return nil
}

func (s *archivistServer) GetSubjects(req *archivist.GetSubjectsRequest, server archivist.Archivist_GetSubjectsServer) error {
	ctx, cancel := context.WithCancel(server.Context())
	defer cancel()
	subjects, err := s.store.GetSubjects(ctx, req)
	if err != nil {
		return err
	}

	for subject := range subjects {
		if err := server.Send(subject); err != nil {
			return err
		}
	}

	return nil
}

type collectorServer struct {
	archivist.UnimplementedCollectorServer

	metadataStore MetadataStorer
	objectStore   ObjectStorer
}

type Storer interface {
	Store(context.Context, string, []byte) error
}

type MetadataStorer interface {
	Storer
	GetBySubjectDigest(context.Context, *archivist.GetBySubjectDigestRequest) (<-chan *archivist.GetBySubjectDigestResponse, error)
	GetSubjects(context.Context, *archivist.GetSubjectsRequest) (<-chan *archivist.GetSubjectsResponse, error)
}

type ObjectStorer interface {
	Storer
	Get(context.Context, *archivist.GetRequest) (io.ReadCloser, error)
}

func NewCollectorServer(metadataStore MetadataStorer, objectStore ObjectStorer) archivist.CollectorServer {
	return &collectorServer{
		objectStore:   objectStore,
		metadataStore: metadataStore,
	}
}

func (s *collectorServer) Store(server archivist.Collector_StoreServer) error {
	ctx := server.Context()
	payload := make([]byte, 0)
	for {
		c, err := server.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		payload = append(payload, c.GetChunk()...)
	}

	// generate gitoid
	gid, err := gitoid.New(bytes.NewBuffer(payload), gitoid.WithContentLength(int64(len(payload))), gitoid.WithSha256())
	if err != nil {
		log.FromContext(ctx).Errorf("failed to generate gitoid: %v", err)
		return err
	}

	if err := s.metadataStore.Store(ctx, gid.String(), payload); err != nil {
		log.FromContext(ctx).Errorf("received error from metadata store: %+v", err)
		return err
	}

	if s.objectStore != nil {
		if err := s.objectStore.Store(ctx, gid.String(), payload); err != nil {
			log.FromContext(ctx).Errorf("received error from object store: %+v", err)
			return err
		}
	}

	return server.SendAndClose(&archivist.StoreResponse{Gitoid: gid.String()})
}

func (s *collectorServer) Get(request *archivist.GetRequest, server archivist.Collector_GetServer) error {
	if s.objectStore == nil {
		return s.UnimplementedCollectorServer.Get(request, server)
	}

	objReader, err := s.objectStore.Get(server.Context(), request)
	if err != nil {
		return err
	}

	defer objReader.Close()
	chunk := &archivist.Chunk{}
	buf := make([]byte, ChunkSize)
	r := bufio.NewReaderSize(objReader, ChunkSize)
	for {
		n, err := io.ReadFull(r, buf)
		if err == io.EOF {
			break
		} else if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		chunk.Chunk = buf[:n]
		if err := server.Send(chunk); err != nil {
			return err
		}
	}

	return nil
}
