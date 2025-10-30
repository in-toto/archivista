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

package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIamCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "PSQL No IAM",
			args:           []string{"iam", "PSQL", "postgres://user:password@host:5432/dbname"},
			expectedOutput: "postgres://user:password@host:5432/dbname\n",
		},
		{
			name:           "PSQL RDS IAM",
			args:           []string{"iam", "PSQL_RDS_IAM", "postgres://user@host:5432/dbname", "--dryrun"},
			expectedOutput: "postgres://user:authtoken@host:5432/dbname\n",
		},
		{
			name:           "MYSQL RDS IAM",
			args:           []string{"iam", "MYSQL_RDS_IAM", "user@tcp(host)/dbname", "--dryrun"},
			expectedOutput: "user:authtoken@tcp(host:3306)/dbname?allowCleartextPasswords=true&parseTime=true&tls=true\n",
		},
		{
			name:          "PSQL RDS IAM Invalid Connection String",
			args:          []string{"iam", "PSQL_RDS_IAM", "http://invalid", "--dryrun"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()

			w.Close()
			os.Stdout = old

			var out bytes.Buffer
			_, copyErr := io.Copy(&out, r)
			assert.NoError(t, copyErr)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, out.String())
			}
		})
	}
}
