// Copyright 2024 The Archivista Contributors
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

// A note: this follows a pattern followed by network service mesh.
// The pattern was copied from the Network Service Mesh Project
// and modified for use here. The original code was published under the
// Apache License V2.

package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/in-toto/archivista/pkg/artifactstore"
	"github.com/in-toto/archivista/pkg/config"
	"github.com/in-toto/archivista/pkg/metadatastorage/sqlstore"
	"github.com/in-toto/archivista/pkg/objectstorage/blobstore"
	"github.com/in-toto/archivista/pkg/objectstorage/filestore"
	"github.com/in-toto/archivista/pkg/publisherstore"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// Service is the interface for the Archivista service
type Service interface {
	Setup() (Server, error)
	GetConfig() *config.Config
	GetFileStoreCh() chan error
	GetSQLStoreCh() chan error
}

// ArchivistaService is the implementation of the Archivista service
type ArchivistaService struct {
	Ctx         context.Context // context for the service
	Cfg         *config.Config  // configuration for the service (if none it uses environment variables)
	fileStoreCh <-chan error
	sqlStoreCh  <-chan error
}

// Setup Archivista Service
func (a *ArchivistaService) Setup() (*Server, error) {
	var (
		level          logrus.Level
		err            error
		sqlStore       *sqlstore.Store
		fileStore      StorerGetter
		publisherStore []publisherstore.Publisher
	)
	serverOpts := make([]Option, 0)

	startTime := time.Now()
	now := time.Now()
	if a.Cfg == nil {

		logrus.Infof("executing: get config from environment (time since start: %s)", time.Since(startTime))

		a.Cfg = new(config.Config)
		if err := a.Cfg.Process(); err != nil {
			logrus.Fatal(err)
		}
		level, err = logrus.ParseLevel(a.Cfg.LogLevel)
		if err != nil {
			logrus.Fatalf("invalid log level %s", a.Cfg.LogLevel)
		}
		logrus.WithField("duration", time.Since(now)).Infof("completed phase: get config from environment")
	} else {
		logrus.Infof("executing: load given config (time since start: %s)", time.Since(startTime))
		level, err = logrus.ParseLevel(a.Cfg.LogLevel)
		if err != nil {
			logrus.Fatalf("invalid log level %s", a.Cfg.LogLevel)
		}
		logrus.WithField("duration", time.Since(now)).Infof("completed phase: load given config")
	}
	logrus.SetLevel(level)

	// ********************************************************************************
	logrus.Infof("executing phase: initializing storage clients (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()
	fileStore, a.fileStoreCh, err = a.initObjectStore()
	if err != nil {
		logrus.Fatalf("could not create object store: %+v", err)
	}
	serverOpts = append(serverOpts, WithObjectStore(fileStore))

	if a.Cfg.EnableSQLStore {
		entClient, err := sqlstore.NewEntClient(
			a.Cfg.SQLStoreBackend,
			a.Cfg.SQLStoreConnectionString,
			sqlstore.ClientWithMaxIdleConns(a.Cfg.SQLStoreMaxIdleConnections),
			sqlstore.ClientWithMaxOpenConns(a.Cfg.SQLStoreMaxOpenConnections),
			sqlstore.ClientWithConnMaxLifetime(a.Cfg.SQLStoreConnectionMaxLifetime),
		)
		if err != nil {
			logrus.Fatalf("could not create ent client: %+v", err)
		}

		// Continue with the existing setup code for the SQLStore
		sqlStore, a.sqlStoreCh, err = sqlstore.New(context.Background(), entClient)
		if err != nil {
			logrus.Fatalf("error initializing new SQLStore: %+v", err)
		}
		serverOpts = append(serverOpts, WithMetadataStore(sqlStore))

		// Add SQL client for ent
		sqlClient := sqlStore.GetClient()
		serverOpts = append(serverOpts, WithEntSqlClient(sqlClient))
	} else {
		sqlStoreChan := make(chan error)
		a.sqlStoreCh = sqlStoreChan
		go func() {
			<-a.Ctx.Done()
			close(sqlStoreChan)
		}()
	}

	// initialize the artifact store
	if a.Cfg.EnableArtifactStore {
		wds, err := artifactstore.New(artifactstore.WithConfigFile(a.Cfg.ArtifactStoreConfig))
		if err != nil {
			logrus.Fatalf("could not create the artifact store: %+v", err)
		}

		serverOpts = append(serverOpts, WithArtifactStore(wds))
	}

	if a.Cfg.Publisher != nil {
		publisherStore = publisherstore.New(a.Cfg)
		serverOpts = append(serverOpts, WithPublishers(publisherStore))
	}
	// Create the Archivista server with all options
	server, err := New(a.Cfg, serverOpts...)
	if err != nil {
		logrus.Fatalf("could not create archivista server: %+v", err)
	}

	logrus.WithField("duration", time.Since(now)).Infof("completed phase: initializing storage clients")
	return &server, nil
}

// GetFileStoreCh returns the file store channel
func (a *ArchivistaService) GetFileStoreCh() <-chan error {
	return a.fileStoreCh
}

// GetSQLStoreCh returns the SQL store channel
func (a *ArchivistaService) GetSQLStoreCh() <-chan error {
	return a.sqlStoreCh
}

func (a *ArchivistaService) initObjectStore() (StorerGetter, <-chan error, error) {
	switch strings.ToUpper(a.Cfg.StorageBackend) {
	case "FILE":
		return filestore.New(a.Ctx, a.Cfg.FileDir, a.Cfg.FileServeOn)

	case "BLOB":
		var creds *credentials.Credentials

		switch a.Cfg.BlobStoreCredentialType {
		case "IAM":
			creds = credentials.NewIAM("")
		case "ACCESS_KEY":
			creds = credentials.NewStaticV4(a.Cfg.BlobStoreAccessKeyId, a.Cfg.BlobStoreSecretAccessKeyId, "")
		default:
			logrus.Fatalf("invalid blob store credential type: %s", a.Cfg.BlobStoreCredentialType)
		}
		return blobstore.New(
			a.Ctx,
			a.Cfg.BlobStoreEndpoint,
			creds,
			a.Cfg.BlobStoreBucketName,
			a.Cfg.BlobStoreUseTLS,
		)

	case "":
		errCh := make(chan error)
		go func() {
			<-a.Ctx.Done()
			close(errCh)
		}()
		return nil, errCh, nil

	default:
		return nil, nil, fmt.Errorf("unknown storage backend: %s", a.Cfg.StorageBackend)
	}
}
