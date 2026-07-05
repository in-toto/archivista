// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in the parent credentialrefresh directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package awsrds

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func TestStoreValidation(t *testing.T) {
	if _, err := NewStore(nil); err == nil {
		t.Fatal("expected an error but didn't get one")
	} else if !errors.Is(err, errMissingConfig) {
		t.Fatalf("expected '%T' but got '%T' instead", errMissingConfig, err)
	}

	testCases := []struct {
		expectedErr error
		conf        Config
		description string
	}{
		{
			description: "missing endpoint",
			conf:        Config{Region: "us-east-1", User: "bar"},
			expectedErr: &errMissingConfigItem{item: "endpoint"},
		},
		{
			description: "missing region",
			conf:        Config{Endpoint: "foo", User: "bar"},
			expectedErr: &errMissingConfigItem{item: "region"},
		},
		{
			description: "missing user",
			conf:        Config{Endpoint: "foo", Region: "us-east-1"},
			expectedErr: &errMissingConfigItem{item: "user"},
		},
		{
			description: "malformed endpoint - no port",
			conf:        Config{Endpoint: "foo", Region: "us-east-1", User: "bar"},
			expectedErr: errMalformedEndpoint,
		},
		{
			description: "malformed endpoint - non-numeric port",
			conf:        Config{Endpoint: "foo:bar", Region: "us-east-1", User: "bar"},
			expectedErr: &url.Error{
				Op:  "parse",
				URL: "http://foo:bar",
				Err: errors.New(`invalid port ":bar" after host`),
			},
		},
		{
			description: "malformed endpoint - missing hostname",
			conf:        Config{Endpoint: "http://:5432", Region: "us-east-1", User: "bar"},
			expectedErr: errMalformedEndpoint,
		},
		{
			description: "missing credentials",
			conf:        Config{Endpoint: "localhost:5432", Region: "us-east-1", User: "bar"},
			expectedErr: errMissingCredentials,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			conf := testCase.conf
			_, err := NewStore(&conf)
			if err == nil {
				t.Fatal("expected an error but didn't get one")

				return
			}

			// If we have a pointer to an error we need to compare error strings
			if reflect.ValueOf(testCase.expectedErr).Kind() == reflect.Pointer &&
				err.Error() != testCase.expectedErr.Error() {
				t.Fatalf("expected '%v' but got '%v' instead", testCase.expectedErr, err)

				return
			}

			if reflect.ValueOf(testCase.expectedErr).Kind() != reflect.Pointer &&
				err != testCase.expectedErr {
				t.Fatalf("expected '%T' but got '%T' instead", testCase.expectedErr, err)
			}
		})
	}

	if _, err := NewStore(&Config{
		Endpoint:    "http://localhost:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: aws.AnonymousCredentials{},
	}); err != nil {
		t.Fatalf("expected no error but got %v instead", err)
	}
}

func TestValidStoreCanGenerateToken(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if creds.GetUsername() == "" {
		t.Fatal("got empty username")
	}

	if creds.GetPassword() == "" {
		t.Fatal("got empty password")
	}
}

func TestStoreErrorsOnUnsignableCredentials(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: aws.AnonymousCredentials{},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := s.Get(ctx); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestStoreCachesCredentials(t *testing.T) {
	s, err := NewStore(&Config{
		Endpoint:    "rdsmysql.cdgmuqiadpid.us-east-1.rds.amazonaws.com:5432",
		Region:      "us-east-1",
		User:        "dbuser",
		Credentials: credentials.NewStaticCredentialsProvider("foo", "bar", "baz"),
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	creds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	username := creds.GetUsername()
	if username == "" {
		t.Fatal("got empty username")
	}

	password := creds.GetPassword()
	if password == "" {
		t.Fatal("got empty password")
	}

	// Second time through we should have everything cached
	cachedCreds, err := s.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if username != cachedCreds.GetUsername() {
		t.Fatalf("expected username to be cached but got %s instead", cachedCreds.GetUsername())
	}
	if password != cachedCreds.GetPassword() {
		t.Fatalf("expected password to be cached but got %s instead", cachedCreds.GetPassword())
	}

	// Refresh should produce credentials without error. The upstream test also asserted
	// that the refreshed password differs from the cached one, which required clock
	// manipulation via the archived bou.ke/monkey library. That assertion is omitted here.
	if _, err := s.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
}
