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
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/in-toto/archivista/pkg/metadatastorage/sqlstore"
	"github.com/spf13/cobra"
)

var iamCmd = &cobra.Command{
	Use:                   "iam [sql backend] [connection string]",
	Short:                 "Converts a connection string to use AWS RDS IAM authentication",
	Long:                  `If the sql backend ends with _RDS_IAM, an IAM authentication token is added to the connection string.`,
	Example:               `archivistactl iam PSQL_RDS_IAM "postgres://user@host:3306/dbname"`,
	Args:                  cobra.ExactArgs(2),
	ValidArgs:             []string{"MYSQL", "MYSQL_RDS_IAM", "PSQL", "PSQL_RDS_IAM"},
	Annotations:           map[string]string{"help:args": "sql backend and connection string"},
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		sqlBackend := strings.ToUpper(args[0])
		connectionString := args[1]
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		if strings.HasSuffix(sqlBackend, "_RDS_IAM") {
			if dryrun {
				sqlstore.AwsConfigAPI = &dryrunConfig{
					cfg: aws.Config{
						Region:      "us-east-1",
						Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
					},
				}
			}
			revisedConnectionString, err := sqlstore.RewriteConnectionStringForIAM(sqlBackend, connectionString, dryrun)
			if err != nil {
				return err
			}
			fmt.Println(revisedConnectionString)
		} else {
			fmt.Println(connectionString)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.Flags().Bool("dryrun", false, "Shows the result using a fake authentication token 'authtoken'")
}

type dryrunConfig struct {
	cfg aws.Config
}

func (m *dryrunConfig) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return m.cfg, nil
}
