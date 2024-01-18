// Copyright 2022 The Witness Contributors
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

package source

import (
	"context"

	"github.com/in-toto/go-witness/archivista"
)

type ArchivistaSource struct {
	client      *archivista.Client
	seenGitoids []string
}

func NewArchvistSource(client *archivista.Client) *ArchivistaSource {
	return &ArchivistaSource{
		client:      client,
		seenGitoids: make([]string, 0),
	}
}

func (s *ArchivistaSource) Search(ctx context.Context, collectionName string, subjectDigests, attestations []string) ([]CollectionEnvelope, error) {
	gitoids, err := s.client.SearchGitoids(ctx, archivista.SearchGitoidVariables{
		CollectionName: collectionName,
		SubjectDigests: subjectDigests,
		Attestations:   attestations,
		ExcludeGitoids: s.seenGitoids,
	})

	if err != nil {
		return []CollectionEnvelope{}, err
	}

	envelopes := make([]CollectionEnvelope, 0, len(gitoids))
	for _, gitoid := range gitoids {
		env, err := s.client.Download(ctx, gitoid)
		if err != nil {
			return envelopes, err
		}

		s.seenGitoids = append(s.seenGitoids, gitoid)
		collectionEnv, err := envelopeToCollectionEnvelope(gitoid, env)
		if err != nil {
			return envelopes, err
		}

		envelopes = append(envelopes, collectionEnv)
	}

	return envelopes, nil
}
