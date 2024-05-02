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
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test Suite: UT Retrieve
type UTRetrieveSuite struct {
	suite.Suite
}

func TestUTRetrieveSuite(t *testing.T) {
	suite.Run(t, new(UTRetrieveSuite))
}

func (ut *UTRetrieveSuite) Test_MissingSubCommand() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"retrieve"})
	err := rootCmd.Execute()
	if err != nil {
		ut.FailNow("Expected: error")
	}
	ut.Contains(output.String(), "archivistactl retrieve")
}

func (ut *UTRetrieveSuite) Test_RetrieveEnvelopeMissingArg() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"retrieve", "envelope"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "accepts 1 arg(s), received 0")
	} else {
		ut.FailNow("Expected: error")
	}
}

func (ut *UTRetrieveSuite) Test_RetrieveEnvelope_NoDB() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"retrieve", "envelope", "test"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "connection refused")
	} else {
		ut.FailNow("Expected: error")
	}
}

func (ut *UTRetrieveSuite) Test_RetrieveSubjectsMissingArg() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"retrieve", "subjects"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "accepts 1 arg(s), received 0")
	} else {
		ut.FailNow("Expected: error")
	}
}

func (ut *UTRetrieveSuite) Test_RetrieveSubjectsNoDB() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"retrieve", "subjects", "test"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "connection refused")
	} else {
		ut.FailNow("Expected: error")
	}
}
