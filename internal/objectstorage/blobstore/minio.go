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

package blobstore

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

type Store struct {
	client   *minio.Client
	bucket   string
	location string
}

// PutBlob stores the attestation blob into the backend store
func (store *Store) PutBlob(idx string, obj []byte) error {
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

// New returns a reader/writer for storing/retrieving attestations
func New(ctx context.Context, endpoint, accessKeyId, secretAccessKeyId, bucketName string, useTLS bool) (*Store, <-chan error, error) {
	errCh := make(chan error)
	go func() {
		<-ctx.Done()
		close(errCh)
	}()

	c, err := minio.NewWithOptions(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKeyId, ""),
		Secure: useTLS,
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

	return &Store{
		client:   c,
		location: loc,
		bucket:   bucketName,
	}, errCh, nil
}

func (s *Store) Get(ctx context.Context, gitoid string) (io.ReadCloser, error) {
	return s.client.GetObjectWithContext(ctx, s.bucket, gitoid, minio.GetObjectOptions{})
}

func (s *Store) Store(ctx context.Context, gitoid string, payload []byte) error {
	return s.PutBlob(gitoid, payload)
}
