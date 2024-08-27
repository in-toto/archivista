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
package publisherstore

import (
	"context"
	"strings"

	"github.com/in-toto/archivista/pkg/config"
	"github.com/sirupsen/logrus"
)

type Publisher interface {
	Publish(ctx context.Context, gitoid string, payload []byte) error
}

func New(config *config.Config) []Publisher {
	var publisherStore []Publisher
	for _, pubType := range config.Publisher {
		pubType = strings.ToUpper(pubType) // Normalize the input
		switch pubType {
		// cases here
		default:
			logrus.Errorf("unsupported publisher type: %s", pubType)
		}
	}
	return publisherStore
}
