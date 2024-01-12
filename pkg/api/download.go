// Copyright 2023 The Witness Contributors
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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/in-toto/go-witness/dsse"
)

func Download(ctx context.Context, baseUrl string, gitoid string) (dsse.Envelope, error) {
	buf := &bytes.Buffer{}
	if err := DownloadWithWriter(ctx, baseUrl, gitoid, buf); err != nil {
		return dsse.Envelope{}, err
	}

	env := dsse.Envelope{}
	dec := json.NewDecoder(buf)
	if err := dec.Decode(&env); err != nil {
		return env, err
	}

	return env, nil
}

func DownloadWithWriter(ctx context.Context, baseUrl, gitoid string, dst io.Writer) error {
	downloadUrl, err := url.JoinPath(baseUrl, "download", gitoid)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(errMsg))
	}

	_, err = io.Copy(dst, resp.Body)
	return err
}
