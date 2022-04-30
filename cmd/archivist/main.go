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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"github.com/testifysec/archivist/internal/config"
	"github.com/testifysec/archivist/internal/server"
	"github.com/testifysec/archivist/internal/storage/badgerstore"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/networkservicemesh/sdk/pkg/tools/debug"
	"github.com/networkservicemesh/sdk/pkg/tools/grpcutils"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"github.com/networkservicemesh/sdk/pkg/tools/log/logruslogger"
	"github.com/sirupsen/logrus"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// ********************************************************************************
	// Setup logging
	// ********************************************************************************
	logrus.SetFormatter(&nested.Formatter{})
	log.EnableTracing(true)
	ctx = log.WithLog(ctx, logruslogger.New(ctx, map[string]interface{}{"cmd": os.Args[0]}))

	// ********************************************************************************
	// Debug self if necessary
	// ********************************************************************************

	if err := debug.Self(); err != nil {
		log.FromContext(ctx).Infof("%s", err)
	}

	startTime := time.Now()

	// enumerating phases

	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 1: get config from environment (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now := time.Now()

	cfg := new(config.Config)
	if err := cfg.Process(); err != nil {
		logrus.Fatal(err)
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("invalid log level %s", cfg.LogLevel)
	}
	logrus.SetLevel(level)

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 1: get config from environment")
	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 2: get spiffe svid (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()
	var source *workloadapi.X509Source
	var svid *x509svid.SVID
	if cfg.EnableSPIFFE == true {
		source, err = workloadapi.NewX509Source(ctx)
		if err != nil {
			logrus.Fatalf("error getting x509 source: %+v", err)
		}
		svid, err = source.GetX509SVID()
		if err != nil {
			logrus.Fatalf("error getting x509 svid: %+v", err)
		}
		logrus.Infof("SVID: %q", svid.ID)
		log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 2: retrieve spiffe svid")
	} else {
		log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 2: SKIPPED")
	}
	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 3: initializing badger (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()

	store, storeCh, err := badgerstore.NewServer(ctx, "")
	if err != nil {
		logrus.Fatalf("error starting badger store: %+v", err)
	}

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 3: initializing badger")
	// ********************************************************************************
	log.FromContext(ctx).Infof("executing phase 4: create and register grpc service (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now = time.Now()

	grpcOptions := make([]grpc.ServerOption, 0)
	if cfg.EnableSPIFFE == true {
		grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewTLS(tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeAny()))))
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	svc := server.NewService(store)
	archivist.RegisterArchivistServer(grpcServer, svc)

	srvErrCh := grpcutils.ListenAndServe(ctx, &cfg.ListenOn, grpcServer)
	exitOnErrCh(ctx, cancel, srvErrCh)

	log.FromContext(ctx).WithField("duration", time.Since(now)).Infof("completed phase 4: create and register grpc server")

	log.FromContext(ctx).Infof("startup complete (time since start: %s)", time.Since(startTime))

	<-ctx.Done()
	<-srvErrCh
	<-storeCh

	//// ********************************************************************************
	log.FromContext(ctx).Infof("exiting, uptime: %v", time.Since(startTime))
	//// ********************************************************************************
}

func exitOnErrCh(ctx context.Context, cancel context.CancelFunc, errCh <-chan error) {
	// If we already have an error, log it and exit
	select {
	case err := <-errCh:
		log.FromContext(ctx).Fatal(err)
	default:
	}
	// Otherwise, wait for an error in the background to log and cancel
	go func(ctx context.Context, errCh <-chan error) {
		err := <-errCh
		if err != nil {
			log.FromContext(ctx).Error(err)
		}
		cancel()
	}(ctx, errCh)
}
