package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
)

const (
	username = "foo"
	password = "bar"
	host     = "localhost"
	port     = 3306
)

type testStore struct {
	Getter    func(ctx context.Context) (Credentials, error)
	Refresher func(ctx context.Context) (Credentials, error)
}

func (vs *testStore) Get(ctx context.Context) (Credentials, error) {
	return vs.Getter(ctx)
}

func (vs *testStore) Refresh(ctx context.Context) (Credentials, error) {
	return vs.Refresher(ctx)
}

type testDriver struct {
	Called         int
	ConnStr        string
	Conn           driver.Conn
	ConnErr        error
	Connector      driver.Connector
	ConnectorError error
}

func (d *testDriver) Open(dsn string) (driver.Conn, error) {
	d.Called++
	d.ConnStr = dsn

	return d.Conn, d.ConnErr
}

func (d *testDriver) OpenConnector(dsn string) (driver.Connector, error) {
	return d.Connector, d.ConnectorError
}

type testFailingDriver struct {
	Called  int
	ConnErr error
}

func (fd *testFailingDriver) Open(dsn string) (driver.Conn, error) {
	fd.Called++
	if fd.Called == 1 {
		return nil, fd.ConnErr
	}

	return nil, nil
}

type testRetryingFailureDriver struct {
	Called    int
	MaxCalled int
	ConnErr   []error
}

func (rfd *testRetryingFailureDriver) Open(dsn string) (driver.Conn, error) {
	if rfd.MaxCalled == 0 {
		rfd.MaxCalled = 1
	}

	rfd.Called++
	if rfd.Called <= rfd.MaxCalled {
		return nil, rfd.ConnErr[rfd.Called-1]
	}

	return nil, nil
}

type testCredential struct {
	Username string
	Password string
}

func (c *testCredential) GetUsername() string {
	return c.Username
}

func (c *testCredential) GetPassword() string {
	return c.Password
}

func TestNewConnectorFailsWithNilConfig(t *testing.T) {
	unregisterAllDrivers()
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    &testDriver{},
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	getFn := func(ctx context.Context) (Credentials, error) {
		return nil, errors.New("error getting creds")
	}

	if _, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", nil); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestNewConnectorWithInvalidDriver(t *testing.T) {
	unregisterAllDrivers()

	getFn := func(ctx context.Context) (Credentials, error) {
		return nil, nil
	}
	if _, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", &Config{
		Host: "",
		Port: 0,
		DB:   "",
		Opts: nil,
	}); err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func TestConnectorErrorsIfStoreGetFailsReturnsNilOrIsInvalid(t *testing.T) {
	unregisterAllDrivers()
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    &testDriver{},
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Host: host,
		Port: port,
		DB:   "test",
	}

	testCases := []struct {
		description string
		getFn       func(ctx context.Context) (Credentials, error)
	}{
		{
			description: "error getting creds",
			getFn: func(ctx context.Context) (Credentials, error) {
				return nil, errors.New("error getting creds")
			},
		},
		{
			description: "nil response",
			getFn: func(ctx context.Context) (Credentials, error) {
				return nil, nil
			},
		},
		{
			description: "empty credentials",
			getFn: func(ctx context.Context) (Credentials, error) {
				return &testCredential{
					Username: "",
					Password: "",
				}, nil
			},
		},
		{
			description: "username missing",
			getFn: func(ctx context.Context) (Credentials, error) {
				return &testCredential{
					Username: "",
					Password: password,
				}, nil
			},
		},
		{
			description: "password missing",
			getFn: func(ctx context.Context) (Credentials, error) {
				return &testCredential{
					Username: username,
					Password: "",
				}, nil
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			c, err := NewConnector(&testStore{
				Getter:    testCase.getFn,
				Refresher: testCase.getFn,
			}, "driver", cfg)
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()

			if _, err = c.Connect(ctx); err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}

func TestConnectorCanUseAlternateFormatter(t *testing.T) {
	unregisterAllDrivers()
	d := &testDriver{}
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	db := "test"

	getFn := func(ctx context.Context) (Credentials, error) {
		return &testCredential{
			Username: username,
			Password: password,
		}, nil
	}
	ctx := context.Background()

	c, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", &Config{
		Host:      host,
		Port:      port,
		DB:        db,
		Formatter: PgKVFormatter,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	dsn := PgKVFormatter(username, password, host, port, db, nil)
	if d.ConnStr != dsn {
		t.Fatalf("expected %s but got %s instead", dsn, d.ConnStr)
	}
}

func TestConnectorRefreshesCredentialsCorrectly(t *testing.T) {
	unregisterAllDrivers()
	d := &testDriver{}
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	getFn := func(ctx context.Context) (Credentials, error) {
		return &testCredential{
			Username: username,
			Password: password,
		}, nil
	}

	c, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", &Config{
		Host: host,
		Port: port,
		DB:   "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	if d.Called != 1 {
		t.Fatalf("expected driver.Open to only have been called once but it was called %d times", d.Called)
	}
}

func TestConnectorFailsToConnectThenReconnects(t *testing.T) {
	unregisterAllDrivers()
	d := &testFailingDriver{
		ConnErr: errors.New(MysqlErrorText),
	}
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	getFn := func(ctx context.Context) (Credentials, error) {
		return &testCredential{
			Username: username,
			Password: password,
		}, nil
	}

	c, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", &Config{
		Host: host,
		Port: port,
		DB:   "test",
		Opts: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	if d.Called != 2 {
		t.Fatalf("expected driver.Open to only have been called twice but it was called %d times", d.Called)
	}
}

func TestConnectorFailsToRefreshOnConnectionFailure(t *testing.T) {
	unregisterAllDrivers()
	d := &testFailingDriver{
		ConnErr: errors.New(MysqlErrorText),
	}
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host: host,
		Port: port,
		DB:   "test",
		Opts: nil,
	}

	refreshCalled := 0

	c, err := NewConnector(&testStore{
		Getter: func(ctx context.Context) (Credentials, error) {
			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
		Refresher: func(ctx context.Context) (Credentials, error) {
			refreshCalled++

			return nil, errors.New("failed to refresh creds")
		},
	}, "driver", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err == nil {
		t.Fatal("expected an error but got nil")
	}

	if d.Called != 1 {
		t.Fatalf("expected driver.Open to only have been called once but it was called %d times", d.Called)
	}

	if refreshCalled != 1 {
		t.Fatalf("expected Refresh func to have been called once but it was called %d times", refreshCalled)
	}
}

func TestConnectorRetriesUntilSuccess(t *testing.T) {
	unregisterAllDrivers()
	mysqlErr := errors.New(MysqlErrorText)
	d := &testRetryingFailureDriver{
		ConnErr:   []error{mysqlErr, mysqlErr, mysqlErr},
		MaxCalled: 3,
	}

	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host:    host,
		Port:    port,
		DB:      "test",
		Opts:    nil,
		Retries: 3,
	}

	refreshCalled := 0

	c, err := NewConnector(&testStore{
		Getter: func(ctx context.Context) (Credentials, error) {
			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
		Refresher: func(ctx context.Context) (Credentials, error) {
			refreshCalled++
			if refreshCalled <= 3 {
				return &testCredential{
					Username: username,
					Password: password,
				}, nil
			}

			return nil, errors.New("failed to refresh creds")
		},
	}, "driver", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err != nil {
		t.Fatal("expected an error but got nil")
	}

	expectedOpenCalls := cfg.Retries + 1
	if d.Called != expectedOpenCalls {
		t.Fatalf("expected driver.Open to only have been called %d time but it was called %d times", expectedOpenCalls, d.Called)
	}

	if refreshCalled != int(cfg.Retries) {
		t.Fatalf("expected Refresh func to have been called %d time but it was called %d times", cfg.Retries, refreshCalled)
	}
}

func TestConnectorRetriesUntilMax(t *testing.T) {
	unregisterAllDrivers()
	var connErr []error
	maxCalled := 5
	for i := 0; i < maxCalled; i++ {
		connErr = append(connErr, errors.New(MysqlErrorText))
	}

	d := &testRetryingFailureDriver{
		ConnErr:   connErr,
		MaxCalled: maxCalled,
	}

	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host:    host,
		Port:    port,
		DB:      "test",
		Opts:    nil,
		Retries: 2,
	}

	refreshCalled := 0

	c, err := NewConnector(&testStore{
		Getter: func(ctx context.Context) (Credentials, error) {
			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
		Refresher: func(ctx context.Context) (Credentials, error) {
			refreshCalled++

			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
	}, "driver", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err == nil {
		t.Fatal("expected an error but got nil")
	}

	expectedOpenCalls := cfg.Retries + 1
	if d.Called != expectedOpenCalls {
		t.Fatalf("expected driver.Open to only have been called %d time but it was called %d times", expectedOpenCalls, d.Called)
	}

	if refreshCalled != cfg.Retries {
		t.Fatalf("expected Refresh func to have been called %d time but it was called %d times", cfg.Retries, refreshCalled)
	}
}

func TestConnectorRetriesUntilNonAuthError(t *testing.T) {
	unregisterAllDrivers()
	var connErr []error
	maxCalled := 5
	for i := 0; i < maxCalled-1; i++ {
		connErr = append(connErr, errors.New(MysqlErrorText))
	}
	connErr = append(connErr, errors.New("another db error error"))

	d := &testRetryingFailureDriver{
		ConnErr:   connErr,
		MaxCalled: maxCalled,
	}

	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host:    host,
		Port:    port,
		DB:      "test",
		Opts:    nil,
		Retries: 5,
	}

	refreshCalled := 0

	c, err := NewConnector(&testStore{
		Getter: func(ctx context.Context) (Credentials, error) {
			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
		Refresher: func(ctx context.Context) (Credentials, error) {
			refreshCalled++

			return &testCredential{
				Username: username,
				Password: password,
			}, nil
		},
	}, "driver", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err == nil {
		t.Fatal("expected an error but got nil")
	}

	if d.Called != cfg.Retries {
		t.Fatalf("expected driver.Open to only have been called %d time but it was called %d times", cfg.Retries, d.Called)
	}

	if refreshCalled != cfg.Retries-1 {
		t.Fatalf("expected Refresh func to have been called %d time but it was called %d times", cfg.Retries-1, refreshCalled)
	}
}

func TestConnectorErrorsIfUnknownDBErrorMessage(t *testing.T) {
	unregisterAllDrivers()
	d := &testFailingDriver{
		ConnErr: errors.New(""),
	}
	if err := Register("driver", func() *Driver {
		return &Driver{
			Driver:    d,
			Formatter: MysqlFormatter,
			AuthError: errorTester(MysqlErrorText),
		}
	}); err != nil {
		t.Fatal(err)
	}

	cfg := Config{
		Host: host,
		Port: port,
		DB:   "test",
		Opts: nil,
	}

	getFn := func(ctx context.Context) (Credentials, error) {
		return &testCredential{
			Username: username,
			Password: password,
		}, nil
	}

	c, err := NewConnector(&testStore{
		Getter:    getFn,
		Refresher: getFn,
	}, "driver", &cfg)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	if _, err := c.Connect(ctx); err == nil {
		t.Fatal("expected error but got nil")
	}

	if d.Called != 1 {
		t.Fatalf("expected driver.Open to only have been called once but it was called %d times", d.Called)
	}
}
