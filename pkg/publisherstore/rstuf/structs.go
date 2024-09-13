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

package rstuf

// Hashes represents the Hashes structure
type Hashes struct {
	Sha256 string `json:"sha256"`
}

// ArtifactInfo represents the ArtifactInfo structure
type ArtifactInfo struct {
	Length int            `json:"length"`
	Hashes Hashes         `json:"hashes"`
	Custom map[string]any `json:"custom,omitempty"`
}

// Artifact represents the Artifact structure
type Artifact struct {
	Path string       `json:"path"`
	Info ArtifactInfo `json:"info"`
}

// ArtifactPayload represents the payload structure
type ArtifactPayload struct {
	Artifacts         []Artifact `json:"artifacts"`
	AddTaskIDToCustom bool       `json:"add_task_id_to_custom"`
	PublishTargets    bool       `json:"publish_targets"`
}

type ArtifactsResponse struct {
	Artifacts  []string `json:"artifacts"`
	TaskId     string   `json:"task_id"`
	LastUpdate string   `json:"last_update"`
	Message    string   `json:"message"`
}

type Response struct {
	Data ArtifactsResponse `json:"data"`
}
