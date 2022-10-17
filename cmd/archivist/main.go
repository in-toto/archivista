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

// A note: this follows a pattern followed by network service mesh.
// The pattern was copied from the Network Service Mesh Project
// and modified for use here. The original code was published under the
// Apache License V2.

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/networkservicemesh/sdk/pkg/tools/debug"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"github.com/networkservicemesh/sdk/pkg/tools/log/logruslogger"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist"
	"github.com/testifysec/archivist/internal/config"
	"github.com/testifysec/archivist/internal/metadatastorage/mysqlstore"
	"github.com/testifysec/archivist/internal/objectstorage/blobstore"
	"github.com/testifysec/archivist/internal/objectstorage/filestore"
	"github.com/testifysec/archivist/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	logrus.SetFormatter(&nested.Formatter{})
	log.EnableTracing(true)
	ctx = log.WithLog(ctx, logruslogger.New(ctx, map[string]interface{}{"cmd": os.Args[0]}))

	if err := debug.Self(); err != nil {
		log.FromContext(ctx).Infof("%s", err)
	}
	startTime := time.Now()

	log.FromContext(ctx).Infof("executing phase 1: get config from environment (time since start: %s)", time.Since(startTime))
	now := time.Now()

	cfg := new(config.Config)
	if err := cfg.Process(); err != nil {
		log.FromContext(ctx).Fatal(err)
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.FromContext(ctx).Fatalf("invalid log level %s", cfg.LogLevel)
	}
	logrus.SetLevel(level)

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 1: get config from environment")

	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 2: initializing storage clients (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()
	fileStore, fileStoreCh, err := initObjectStore(ctx, cfg)
	if err != nil {
		log.FromContext(ctx).Fatalf("error initializing storage clients: %+v", err)
	}

	mysqlStore, mysqlStoreCh, err := mysqlstore.New(ctx, cfg.SQLStoreConnectionString)
	if err != nil {
		log.FromContext(ctx).Fatalf("error initializing mysql client: %+v", err)
	}

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 3: initializing storage clients")

	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 3: create and register http service (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()
	server := server.New(mysqlStore, fileStore)
	router := mux.NewRouter()
	router.HandleFunc("/download/{gitoid}", server.GetHandler)
	router.HandleFunc("/upload", server.StoreHandler)

	if cfg.EnableGraphql {
		client := mysqlStore.GetClient()
		srv := handler.NewDefaultServer(archivist.NewSchema(client))
		srv.Use(entgql.Transactioner{TxOpener: client})
		router.Handle("/query", srv)
		if cfg.GraphqlWebClientEnable {
			router.Handle("/",
				playground.Handler("Archivist", "/query"),
			)
		}
	}

	listenAddress := cfg.ListenOn
	listenAddress = strings.ToLower(strings.TrimSpace(listenAddress))
	proto := ""
	if strings.HasPrefix(listenAddress, "tcp://") {
		proto = "tcp"
		listenAddress = listenAddress[len("tcp://"):]
	} else if strings.HasPrefix(listenAddress, "unix://") {
		proto = "unix"
		listenAddress = listenAddress[len("unix://"):]
	}

	listener, err := net.Listen(proto, listenAddress)
	if err != nil {
		log.FromContext(ctx).Fatalf("unable to start http listener: %+v", err)
	}

	go func() {
		if err := http.Serve(listener, handlers.CORS(
			handlers.AllowedOrigins(cfg.CORSAllowOrigins),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}),
		)(router)); err != nil {
			log.FromContext(ctx).Fatalf("unable to start http server: %+v", err)
		}
	}()

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 5: create and register http service")
	log.FromContext(ctx).Infof("startup complete (time since start: %s)", time.Since(startTime))

	<-ctx.Done()
	<-fileStoreCh
	<-mysqlStoreCh

	log.FromContext(ctx).Infof("exiting, uptime: %v", time.Since(startTime))
}

func initObjectStore(ctx context.Context, cfg *config.Config) (server.StorerGetter, <-chan error, error) {
	switch strings.ToUpper(cfg.StorageBackend) {
	case "FILE":
		return filestore.New(ctx, cfg.FileDir, cfg.FileServeOn)

	case "BLOB":
		return blobstore.New(
			ctx,
			cfg.BlobStoreEndpoint,
			cfg.BlobStoreAccessKeyId,
			cfg.BlobStoreSecretAccessKeyId,
			cfg.BlobStoreBucketName,
			cfg.BlobStoreUseTLS,
		)

	case "":
		errCh := make(chan error)
		go func() {
			<-ctx.Done()
			close(errCh)
		}()
		return nil, errCh, nil

	default:
		return nil, nil, fmt.Errorf("unknown storage backend: %s", cfg.StorageBackend)
	}
}
