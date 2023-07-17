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

package witness

// all of the following imports are here so that each of the package's init functions run appropriately
import (
	// attestors
	_ "github.com/testifysec/go-witness/attestation/aws-iid"
	_ "github.com/testifysec/go-witness/attestation/commandrun"
	_ "github.com/testifysec/go-witness/attestation/environment"
	_ "github.com/testifysec/go-witness/attestation/gcp-iit"
	_ "github.com/testifysec/go-witness/attestation/git"
	_ "github.com/testifysec/go-witness/attestation/github"
	_ "github.com/testifysec/go-witness/attestation/gitlab"
	_ "github.com/testifysec/go-witness/attestation/jwt"
	_ "github.com/testifysec/go-witness/attestation/maven"
	_ "github.com/testifysec/go-witness/attestation/oci"
	_ "github.com/testifysec/go-witness/attestation/sarif"

	// signer providers
	_ "github.com/testifysec/go-witness/signer/file"
	_ "github.com/testifysec/go-witness/signer/fulcio"
	_ "github.com/testifysec/go-witness/signer/spiffe"
)
