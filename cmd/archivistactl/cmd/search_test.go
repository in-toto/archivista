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

// Test Suite: UT Search
type UTSearchSuite struct {
	suite.Suite
}

func TestUTSearchSuite(t *testing.T) {
	suite.Run(t, new(UTSearchSuite))
}

func (ut *UTSearchSuite) Test_SearchMissingArgs() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"search"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "expected exactly 1 argument")
	} else {
		ut.FailNow("Expected: error")
	}
}

func (ut *UTSearchSuite) Test_SearchInvalidDigestString() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"search", "invalidDigest"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "invalid digest string. expected algorithm:digest")
	} else {
		ut.FailNow("Expected: error")
	}
}

func (ut *UTSearchSuite) Test_NoDB() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"search", "sha256:test"})
	err := rootCmd.Execute()
	if err != nil {
		ut.ErrorContains(err, "connection refused")
	} else {
		ut.FailNow("Expected: error")
	}
}
