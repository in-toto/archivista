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

package format

var (
	// handlers is a slice of registered format handlers
	// Handlers are checked in order, first match wins
	handlers []Handler
)

// GetHandler finds the first handler that can process the given data
func GetHandler(obj []byte) (Handler, bool) {
	for _, h := range handlers {
		if h.Detect(obj) {
			return h, true
		}
	}
	return nil, false
}
