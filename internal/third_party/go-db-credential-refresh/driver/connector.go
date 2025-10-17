package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"sync"
)

// Config is a struct that holds non-credential database configuration.
type Config struct {
	Opts      map[string]string
	Formatter Formatter
	Host      string
	DB        string
	Port      int
	Retries   int
}

var (
	ErrConfigRequired   = errors.New("config is required")
	ErrNoNilCredentials = errors.New("store cannot return nil credentials")
	ErrMissingUsername  = errors.New("missing username")
	ErrMissingPassword  = errors.New("missing password")
)

// NewConnector creates a new connector from a store.
func NewConnector(s Store, driverName string, cfg *Config) (*Connector, error) {
	if cfg == nil {
		return nil, ErrConfigRequired
	}

	d, err := CreateDriver(driverName)
	if err != nil {
		return nil, err
	}

	// Allow caller to override formatter. This makes it easier to use different DSN
	// formats in cases where a default formatter might be difficult to use.
	if cfg.Formatter != nil {
		d.Formatter = cfg.Formatter
	}

	// 0 retries means that it should try once, retry, then don't attempt any more retries
	if cfg.Retries <= 0 {
		cfg.Retries = 1
	}

	return &Connector{
		store:      s,
		cfg:        cfg,
		driver:     d.Driver,
		errHandler: d.AuthError,
		formatter:  d.Formatter,
		mu:         sync.Mutex{},
	}, nil
}

// Connector represents a driver in a fixed configuration.
type Connector struct {
	store      Store
	cfg        *Config
	driver     driver.Driver
	errHandler AuthError
	formatter  Formatter
	mu         sync.Mutex
}

// Connect implements driver.Connector interface.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	creds, err := c.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	if creds == nil {
		return nil, ErrNoNilCredentials
	}

	username := creds.GetUsername()
	password := creds.GetPassword()

	if username == "" {
		return nil, ErrMissingUsername
	}

	if password == "" {
		return nil, ErrMissingPassword
	}

	connStr := c.formatter(username, password, c.cfg.Host, c.cfg.Port, c.cfg.DB, c.cfg.Opts)

	conn, err := c.driver.Open(connStr)
	if err == nil {
		return conn, nil
	}

	if !c.errHandler(err) {
		return nil, err
	}

	for i := 0; i < c.cfg.Retries; i++ {
		creds, err = c.store.Refresh(ctx)
		if err != nil {
			return nil, err
		}

		connStr = c.formatter(
			creds.GetUsername(),
			creds.GetPassword(),
			c.cfg.Host,
			c.cfg.Port,
			c.cfg.DB,
			c.cfg.Opts,
		)

		conn, err = c.driver.Open(connStr)
		if err == nil {
			return conn, nil
		}

		// Bail if we get an error that we can't handle with new creds
		if !c.errHandler(err) {
			return nil, err
		}
	}

	// If we've exhausted our retries we'll just return the last error
	return nil, err
}

// Driver implements driver.Connector interface.
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
