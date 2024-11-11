// Copyright 2022-2024 The Archivista Contributors
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
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/handlers"
	"github.com/in-toto/archivista/pkg/server"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&nested.Formatter{})
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	startTime := time.Now()

	archivistaService := &server.ArchivistaService{Ctx: ctx, Cfg: nil}

	server, err := archivistaService.Setup()
	if err != nil {
		logrus.Fatalf("unable to setup archivista service: %+v", err)
	}
	// ********************************************************************************
	logrus.Infof("executing phase: create and register http service (time since start: %s)", time.Since(startTime))
	// ********************************************************************************
	now := time.Now()

	listenAddress := archivistaService.Cfg.ListenOn
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
		logrus.Fatalf("unable to start http listener: %+v", err)
	}
	srv := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedOrigins(archivistaService.Cfg.CORSAllowOrigins),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}),
		)(server.Router()),
		ReadTimeout:  time.Duration(archivistaService.Cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(archivistaService.Cfg.WriteTimeout) * time.Second,
	}

	go func() {
		if archivistaService.Cfg.EnableTLS {
			if err := srv.ListenAndServeTLS(archivistaService.Cfg.TLSCert, archivistaService.Cfg.TLSKey); err != nil {
				logrus.Fatalf("unable to start http serveR: %+v", err)
			}
		} else {
			if err := srv.Serve(listener); err != nil {
				logrus.Fatalf("unable to start http server: %+v", err)
			}
		}
	}()

	logrus.WithField("duration", time.Since(now)).Infof("completed phase: create and register http service")
	logrus.Infof("startup complete (time since start: %s)", time.Since(startTime))

	<-ctx.Done()
	<-archivistaService.GetFileStoreCh()
	<-archivistaService.GetSQLStoreCh()

	logrus.Infof("exiting, uptime: %v", time.Since(startTime))
}
