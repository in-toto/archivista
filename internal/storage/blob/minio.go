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

package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/git-bom/gitbom-go"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Indexer calculates the index reference for an input blob,
// and gets/puts blobs at that index in and out of the backing
// blob storage.
type Indexer interface {
	GetRef(obj []byte) (string, error)
	GetBlob(idx string) ([]byte, error)
	PutBlob(idx string, obj []byte) error
}

type UnifiedStorage interface {
	archivist.CollectorServer
	archivist.StorageServer
}

type attestationBlobStore struct {
	archivist.UnimplementedCollectorServer
	archivist.UnimplementedStorageServer

	client   *minio.Client
	bucket   string
	location string
}

// GetRef calculates the index reference for a given object
func (store *attestationBlobStore) GetRef(obj []byte) (string, error) {
	gb := gitbom.NewSha256GitBom()
	if err := gb.AddReference(obj, nil); err != nil {
		return "", err
	}
	return gb.Identity(), nil
}

// GetBlob retrieves an attesation from the backend store
func (store *attestationBlobStore) GetBlob(idx string) ([]byte, error) {
	opt := minio.GetObjectOptions{}
	chunkSize := 8 * 1024
	buf := make([]byte, chunkSize)
	outBuf := bytes.NewBuffer(buf)

	obj, err := store.client.GetObject(store.bucket, idx, opt)
	if err != nil {
		return buf, err
	}

	var n int64
	for {
		_ = opt.SetRange(n, n+int64(chunkSize)-1)
		readBytes, err := outBuf.ReadFrom(obj)
		if err == nil {
			return outBuf.Bytes(), nil
		}
		if err != nil {
			if err == io.EOF {
				_, err = outBuf.ReadFrom(bytes.NewReader(buf))
				break
			}
		}

		n += readBytes
		_, err = outBuf.ReadFrom(bytes.NewReader(buf))
		if err != nil {
			return buf, fmt.Errorf("failed to chunk blob: %v", err)
		}
	}
	return []byte{}, fmt.Errorf("failed to read out object: %v", err)
}

// PutBlob stores the attestation blob into the backend store
func (store *attestationBlobStore) PutBlob(idx string, obj []byte) error {
	opt := minio.PutObjectOptions{}
	size := int64(len(obj))
	n, err := store.client.PutObject(store.bucket, idx, bytes.NewReader(obj), size, opt)
	if err != nil {
		return fmt.Errorf("failed to put blob: %v", err)
	} else if n != size {
		return fmt.Errorf("failed to upload full blob: size %d != uploaded size %d", size, n)
	}
	return nil
}

// NewMinioClient returns a reader/writer for storing/retrieving attestations
func NewMinioClient(ctx context.Context, endpoint, accessKeyId, secretAccessKeyId, bucketName string, useSSL bool) (UnifiedStorage, <-chan error, error) {
	errCh := make(chan error)
	go func() {
		<-ctx.Done()
		close(errCh)
	}()

	c, err := minio.NewWithOptions(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKeyId, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errCh, err
	}

	exists, err := c.BucketExists(bucketName)
	if !exists || err != nil {
		return nil, errCh, fmt.Errorf("failed to find bucket exists: %v", err)
	}

	loc, err := c.GetBucketLocation(bucketName)
	if err != nil {
		return nil, errCh, err
	}

	return &attestationBlobStore{
		client:   c,
		location: loc,
		bucket:   bucketName,
	}, errCh, nil
}

func (s *attestationBlobStore) Get(ctx context.Context, req *archivist.GetRequest) (*archivist.GetResponse, error) {
	obj, err := s.GetBlob(req.GetGitoid())
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to retrieve object: %+v", err)
		return nil, err
	}

	return &archivist.GetResponse{
		Gitoid: req.GetGitoid(),
		Object: string(obj),
	}, nil
}

func (s *attestationBlobStore) Store(ctx context.Context, req *archivist.StoreRequest) (*emptypb.Empty, error) {
	idx, err := s.GetRef([]byte(req.GetObject()))
	if err != nil {
		logrus.WithContext(ctx).Errorf("gitbom tag generation failed: %+v", err)
		return nil, err
	}

	if err := s.PutBlob(idx, []byte(req.GetObject())); err != nil {
		logrus.WithContext(ctx).Errorf("failed to store object: %+v", err)
	}

	return &emptypb.Empty{}, nil
}
