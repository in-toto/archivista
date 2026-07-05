// Copyright 2024 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dapr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/in-toto/archivista/pkg/config"
	"github.com/sirupsen/logrus"
)

type DaprHttp struct {
	Client              *http.Client
	Host                string
	HttpPort            string
	PubsubComponentName string
	PubsubTopic         string
	Url                 string
}

type daprPayload struct {
	Gitoid  string
	Payload []byte
}

type Publisher interface {
	Publish(ctx context.Context, gitoid string, payload []byte) error
}

func (d *DaprHttp) Publish(ctx context.Context, gitoid string, payload []byte) error {
	if d.Client == nil {
		d.Client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	if d.Url == "" {
		d.Url = d.Host + ":" + d.HttpPort +
			"/v1.0/publish/" + d.PubsubComponentName + "/" + d.PubsubTopic
	}

	dp := daprPayload{
		Gitoid:  gitoid,
		Payload: payload,
	}
	// Marshal the message to JSON
	msgBytes, err := json.Marshal(dp)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	res, err := d.Client.Post(d.Url, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		logrus.Printf("failed to publish message: %s", res.Body)
		return fmt.Errorf("failed to publish message: %s", res.Body)
	}
	defer res.Body.Close()

	return nil
}

func NewPublisher(config *config.Config) Publisher {
	daprPublisher := &DaprHttp{
		Host:                config.PublisherDaprHost,
		HttpPort:            config.PublisherDaprPort,
		PubsubComponentName: config.PublisherDaprComponentName,
		PubsubTopic:         config.PublisherDaprTopic,
		Url:                 config.PublisherDaprURL,
	}
	return daprPublisher
}
