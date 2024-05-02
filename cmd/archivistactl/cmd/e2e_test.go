// Copyright 2023 The Archivista Contributors
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

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test Suite: E2E
type E2EStoreSuite struct {
	suite.Suite
}

func TestE2EStoreSuite(t *testing.T) {
	suite.Run(t, new(E2EStoreSuite))
}

// TearDown the container deployment
func (e2e *E2EStoreSuite) TearDownTest() {
	e2e.T().Log("Stopping up containers")
	path, _ := os.Getwd()
	fmt.Print(path)
	cmd := exec.Command("bash", "../../../test/deploy-services.sh", "stop")
	err := cmd.Start()
	if err != nil {
		e2e.FailNow(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		e2e.FailNow(err.Error())
	}
}

// Run the E2E tests
func (e2e *E2EStoreSuite) Test_E2E() {

	// Define tests for supported dbs
	testDBCases := []string{"mysql", "pgsql"}

	// Call the script to deploy the containers for the test cases
	for _, testDB := range testDBCases {
		cmd := exec.Command("bash", "../../../test/deploy-services.sh", "start-"+testDB)
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Start()
		if err != nil {
			e2e.FailNow(err.Error())
		}
		e2e.T().Log("Starting services using DB: " + testDB)
		err = cmd.Wait()
		if err != nil {
			e2e.T().Log(out.String())
			e2e.FailNow(err.Error())
		}

		// define test cases struct
		type testCases struct {
			name                string
			attestation         string // files are stored in test/
			sha256              string
			expectedStore       string
			gitoidStore         string // this value is added during `activistactl store command`
			expectedSearch      string
			expectedError       string
			expectedRetrieveSub string
		}

		// test cases
		testTable := []testCases{
			{
				name:                "valid build attestation",
				attestation:         "../../../test/build.attestation.json",
				sha256:              "423da4cff198bbffbe3220ed9510d32ba96698e4b1f654552521d1f541abb6dc",
				expectedStore:       "stored with gitoid",
				expectedSearch:      "Collection name: build",
				expectedRetrieveSub: "Name: https://witness.dev/attestations/product/v0.1/file:testapp",
			},
			{
				name:                "valid package attestation",
				sha256:              "10cbf0f3d870934921276f669ab707983113f929784d877f1192f43c581f2070",
				attestation:         "../../../test/package.attestation.json",
				expectedStore:       "stored with gitoid",
				expectedSearch:      "Collection name: package",
				expectedRetrieveSub: "Name: https://witness.dev/attestations/git/v0.1/commithash:be20100af602c780deeef50c54f5338662ce917c",
			},
			{
				name:           "duplicated package attestation",
				sha256:         "10cbf0f3d870934921276f669ab707983113f929784d877f1192f43c581f2070",
				attestation:    "../../../test/package.attestation.json",
				expectedStore:  "",
				expectedSearch: "Collection name: package",
				expectedError:  "uplicate",
			},
			{
				name:                "fail attestation",
				attestation:         "../../../test/fail.attestation.json",
				sha256:              "5e8c57df8ae58fe9a29b29f9993e2fc3b25bd75eb2754f353880bad4b9ebfdb3",
				expectedStore:       "stored with gitoid",
				expectedSearch:      "",
				expectedRetrieveSub: "Name: https://witness.dev/attestations/git/v0.1/parenthash:aa35c1f4b1d41c87e139c2d333f09117fd0daf4f",
			},
			{
				name:           "invalid payload attestation",
				attestation:    "../../../test/invalid_payload.attestation.json",
				sha256:         "5e8c57df8ae58fe9a29b29f9993e2fc3b25bd75eb2754f353880bad4b9ebfdb3",
				expectedStore:  "stored with gitoid",
				expectedSearch: "",
				expectedError:  "value is less than the required length",
			},
			{
				name:          "nonexistent payload file",
				attestation:   "../../../test/missing.attestation.json",
				expectedError: "no such file or directory",
			},
		}
		for _, test := range testTable {
			// test `archivistactl store`
			e2e.T().Log("Test `archivistactl store` " + test.name)
			storeOutput := bytes.NewBufferString("")
			rootCmd.SetOut(storeOutput)
			rootCmd.SetErr(storeOutput)
			rootCmd.SetArgs([]string{"store", test.attestation})
			err := rootCmd.Execute()
			if err != nil {
				// if return error assert if is expected error from test case
				e2e.ErrorContains(err, test.expectedError)
			} else { // assert the expected responses
				storeActual := storeOutput.String()
				e2e.Contains(storeActual, test.expectedStore)
				test.gitoidStore = strings.Split(storeActual, "stored with gitoid ")[1]
				test.gitoidStore = strings.TrimSuffix(test.gitoidStore, "\n")
			}

			// test `archivistactl search`
			e2e.T().Log("Test `archivistactl search`" + test.name)
			searchOutput := bytes.NewBufferString("")
			rootCmd.SetOut(searchOutput)
			rootCmd.SetErr(searchOutput)
			rootCmd.SetArgs([]string{"search", "sha256:" + test.sha256})
			err = rootCmd.Execute()
			if err != nil {
				e2e.FailNow(err.Error())
			}
			searchActual := searchOutput.String()
			e2e.Contains(searchActual, test.expectedSearch)

			if test.expectedRetrieveSub != "" {
				// test `archivistactl retrieve subjects`
				e2e.T().Log("Test `archivistactl retrieve subjects` " + test.name)
				subjectsOutput := bytes.NewBufferString("")
				rootCmd.SetOut(subjectsOutput)
				rootCmd.SetErr(subjectsOutput)
				rootCmd.SetArgs([]string{"retrieve", "subjects", test.gitoidStore})
				err = rootCmd.Execute()
				if err != nil {
					e2e.FailNow(err.Error())
				}
				subjectsActual := subjectsOutput.String()
				e2e.Contains(subjectsActual, test.expectedRetrieveSub)
				if test.name == "fail attestation" {
					e2e.NotContains(subjectsActual, "sha256:"+test.sha256)
				} else {
					e2e.Contains(subjectsActual, "sha256:"+test.sha256)
				}
			}
			if test.expectedError == "" {
				tempDir := os.TempDir()
				// test `archivistactl retrieve envelope`
				e2e.T().Log("Test `archivistactl retrieve envelope` " + test.name)
				envelopeOutput := bytes.NewBufferString("")
				rootCmd.SetOut(envelopeOutput)
				rootCmd.SetErr(envelopeOutput)
				rootCmd.SetArgs([]string{"retrieve", "envelope", test.gitoidStore, "-o", path.Join(tempDir, test.gitoidStore)})
				err = rootCmd.Execute()
				if err != nil {
					e2e.FailNow(err.Error())
				}
				// compares file attestation with the retrieved attestation
				fileAtt, err := os.ReadFile(test.attestation)
				if err != nil {
					e2e.FailNow(err.Error())
				}
				fileSaved, err := os.ReadFile(path.Join(tempDir, test.gitoidStore))
				if err != nil {
					e2e.FailNow(err.Error())
				}
				if err != nil {
					e2e.FailNow(err.Error())
				}
				e2e.True(bytes.Equal(fileAtt, fileSaved))
			}
		}
	}

}
