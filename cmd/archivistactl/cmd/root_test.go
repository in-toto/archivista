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

// Test Suite: UT Root
type UTRootSuite struct {
	suite.Suite
}

func TestUTRootSuite(t *testing.T) {
	suite.Run(t, new(UTRootSuite))
}

func (ut *UTRootSuite) Test_Root() {
	output := bytes.NewBufferString("")
	rootCmd.SetOut(output)
	rootCmd.SetErr(output)
	rootCmd.SetArgs([]string{"help"})
	err := rootCmd.Execute()
	if err != nil {
		ut.FailNow(err.Error())
	}
	actual := output.String()
	ut.Contains(actual, "A utility to interact with an archivista server")

}
