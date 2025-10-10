// Copyright 2025 The Archivista Contributors
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

package driverstores

import (
	"context"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/davepgreene/go-db-credential-refresh/store"
)

type Static struct {
	credentials store.Credential
}

func NewStaticStore(username, password string) *Static {
	return &Static{
		credentials: store.Credential{
			Username: username,
			Password: password,
		},
	}
}

func (s Static) Get(ctx context.Context) (driver.Credentials, error) {
	return s.credentials, nil
}

func (s Static) Refresh(ctx context.Context) (driver.Credentials, error) {
	return s.credentials, nil
}
