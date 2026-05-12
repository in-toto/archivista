// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in this directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package driver

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFormatters(t *testing.T) {
	username := "foo"
	password := "bar"
	host := "localhost"

	testCases := []struct {
		name        string
		port        int
		db          string
		opts        map[string]string
		formatter   Formatter
		expectedDsn string
		parseAsURL  bool
	}{
		{
			name: "mysql - with opts",
			port: 3306,
			db:   "test",
			opts: map[string]string{
				"maxAllowedPacket": "8203",
				"tcpKeepAlive":     "true",
			},
			formatter:   MysqlFormatter,
			expectedDsn: "foo:bar@tcp(localhost:3306)/test",
			parseAsURL:  true,
		},
		{
			name:        "mysql - no opts",
			port:        3306,
			formatter:   MysqlFormatter,
			expectedDsn: "foo:bar@tcp(localhost:3306)/",
			parseAsURL:  true,
		},
		{
			name: "pg - with opts",
			port: 5432,
			db:   "test",
			opts: map[string]string{
				"sslmode":         "disable",
				"connect_timeout": "10",
			},
			formatter:   PgFormatter,
			expectedDsn: "postgres://foo:bar@localhost:5432/test",
			parseAsURL:  true,
		},
		{
			name:        "pg - no opts",
			port:        5432,
			formatter:   PgFormatter,
			expectedDsn: "postgres://foo:bar@localhost:5432",
			parseAsURL:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dsn := testCase.formatter(username, password, host, testCase.port, testCase.db, testCase.opts)

			if !testCase.parseAsURL {
				if dsn != testCase.expectedDsn {
					t.Fatalf("expected %s but got %s", testCase.expectedDsn, dsn)
				}

				return
			}
			// Because url.Values is a map[string][]string we get randomized ordering before encoding so we
			// can't compare query strings. Instead we extract the generated DSN's query params to a map and
			// compare with reflect.DeepEqual.
			if testCase.opts != nil {
				u, err := url.Parse(dsn)
				if err != nil {
					t.Fatal(err)
				}

				params := make(map[string]string)
				for k, v := range u.Query() {
					params[k] = v[0]
				}

				if !reflect.DeepEqual(params, testCase.opts) {
					t.Fatalf("expected params %v but got %v", testCase.opts, params)
				}

				u.RawQuery = ""
				dsn = u.String()
			}
			if dsn != testCase.expectedDsn {
				t.Fatalf("expected %s but got %s", testCase.expectedDsn, dsn)
			}
		})
	}
}
