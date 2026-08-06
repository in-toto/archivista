// Copyright 2024 The Archivista Contributors
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
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/edwarnicke/gitoid"
	"github.com/in-toto/archivista/ent"
	"github.com/in-toto/archivista/pkg/metadatastorage/sqlstore"
	"github.com/in-toto/go-witness/dsse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	importCmd = &cobra.Command{
		Use:          "import",
		Short:        "import dsses to the Archivista DB server",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(2),
		RunE:         importDsse,
	}
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().StringP("from-dir", "", "", "Directory to import from. Example: /path/to/directory")
	importCmd.PersistentFlags().StringP("db-uri", "", "", "Database URI to import to. Supported schemes: mysql, psql. Example: mysql://user:password@localhost:3306/testify")
	importCmd.PersistentFlags().IntP("max-concurrent", "", 3, "Maximum number of concurrent imports.")
	err := importCmd.MarkPersistentFlagRequired("db-uri")
	cobra.CheckErr(err)
}

func walkDir(dir string) (<-chan string, <-chan error) {
	ch := make(chan string)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ch <- path
			}
			return nil
		})

		if err != nil {
			errCh <- err
		}
	}()

	return ch, errCh
}

func dbClient(dbURI string) (*ent.Client, error) {
	// Verify that we can connect to the DB
	var (
		scheme string
		uri    string
	)

	purl, err := url.Parse(dbURI)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(purl.Scheme) {
	case "MYSQL":
		scheme = "MYSQL"
		uri = purl.User.String() + "@" + purl.Host + purl.Path
	case "PSQL":
		scheme = "PSQL"
		uri = dbURI
	default:
		return nil, fmt.Errorf("unsupported database scheme %s", purl.Scheme)
	}

	// TODO:
	// Define MaxIdleConns, MaxOpenConns, ConnMaxLifetime as custom parameters
	entClient, err := sqlstore.NewEntClient(
		scheme,
		uri,
	)
	if err != nil {
		return nil, err
	}

	return entClient, nil
}

func importFile(path string, sqlStore *sqlstore.Store, maxConcurrent int) {
	fpaths, _ := walkDir(path)
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent) // Buffered channel acting as a semaphore

	for fpath := range fpaths {
		sem <- struct{}{} // Acquire a token before starting a new Goroutine
		wg.Add(1)
		go func(fpath string) {
			defer wg.Done()
			defer func() { <-sem }() // Release the token when done

			fmt.Println("\nImporting file:", fpath)
			file, err := os.ReadFile(fpath)
			if err != nil {
				fmt.Println("Skipping  file: "+fpath+" cannot read file", fpath)
				return
			}

			envelope := &dsse.Envelope{}
			if err := json.Unmarshal(file, envelope); err != nil {
				fmt.Printf("Skipping file: %s cannot open %s as DSSE Envelope\n", fpath, fpath)
				return
			}
			if envelope.PayloadType != "" {
				ngitoid, err := gitoid.New(bytes.NewReader(file), gitoid.WithContentLength(int64(len(file))), gitoid.WithSha256())
				if err != nil {
					fmt.Println("Skipping  file: "+fpath+" cannot generate valid GitOID", fpath)
					return
				}
				err = sqlStore.Store(context.Background(), ngitoid.String(), file)
				if err != nil {
					// if failed due to duplicate entry, skip
					if strings.Contains(err.Error(), "Duplicate entry") {
						fmt.Println("Skipping  file: " + fpath + " cannot store duplicated entry")
					} else {
						fmt.Println("Skipping  file: "+fpath+" failed to import.", err)
					}
					return
				}
			}
			fmt.Println("Successfully imported", fpath)
		}(fpath)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
}

func importDsse(ccmd *cobra.Command, args []string) error {
	logrus.SetLevel(logrus.FatalLevel) // Set the log level to Fatal to suppress lower-level logs

	sourceDir, err := ccmd.Flags().GetString("from-dir")
	if err != nil {
		return err
	}

	dbURI, err := ccmd.Flags().GetString("db-uri")
	if err != nil {
		return err
	}

	ec, err := dbClient(dbURI)
	if err != nil {
		return err
	}

	sqlStore, _, err := sqlstore.New(context.Background(), ec)
	if err != nil {
		return err
	}

	max, err := ccmd.Flags().GetInt("max-concurrent")
	if err != nil {
		fmt.Println("Failed to get max-concurrent flag", err)
	}
	fmt.Print("\nImporting DSSes from folder", sourceDir, " to the database server")
	fmt.Print("\nMax concurrent imports: ", max, "\n\n")
	importFile(sourceDir, sqlStore, max)

	return nil
}
