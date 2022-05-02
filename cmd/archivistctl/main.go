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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/testifysec/witness/pkg/dsse"
	"io/ioutil"
	"os"

	"github.com/testifysec/archivist-api/pkg/api/archivist"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Fatalln("error")
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		logrus.Fatalf("unable to read file %+v", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		logrus.Fatalf("unable to grpc dial: %+v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logrus.Errorf("unable to close connection: %+v", err)
		}
	}()

	archivistClient := archivist.NewArchivistClient(conn)
	collectorClient := archivist.NewCollectorClient(conn)

	// check if valid

	obj := &dsse.Envelope{}
	err = json.Unmarshal(file, obj)
	if err != nil {
		logrus.Fatalln("could not unmarshal input: ", err)
	}

	if len(obj.Payload) == 0 || obj.PayloadType == "" || len(obj.Signatures) == 0 {
		logrus.Fatalln("obj is not DSSE %d %d %d", len(obj.Payload), len(obj.PayloadType), len(obj.Signatures))
	}

	_, err = collectorClient.Store(context.Background(), &archivist.StoreRequest{
		Object: string(file),
	})
	if err != nil {
		logrus.Errorf("unable to store object: %+v", err)
	}

	resp, err := archivistClient.GetBySubjectDigest(context.Background(), &archivist.GetBySubjectDigestRequest{Algorithm: "sha256", Value: "3ff8d62302fe86d3f2d637843c696ab4d065f9e1cc873121806cf9467e546e4f"})
	if err != nil {
		logrus.Fatalf("unable to retrieve object: %+v", err)
	}

	for i, v := range resp.GetObject() {
		fmt.Printf("%d: %s\n", i, v)
	}
}
