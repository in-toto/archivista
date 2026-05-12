// Vendored from https://github.com/davepgreene/go-db-credential-refresh (MIT).
// See LICENSE in this directory.
//
// Copyright (c) 2022-2024 Dave Greene
// Copyright (c) 2026 The Archivista Contributors

package driver

import (
	"database/sql/driver"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/stdlib"
)

// Driver carries information along with a database/sql/driver required for creating a Connector.
type Driver struct {
	Driver    driver.Driver
	Formatter Formatter
	AuthError AuthError
}

type factory func() *Driver

type errFactoryAlreadyRegistered struct {
	name string
}

func (e errFactoryAlreadyRegistered) Error() string {
	return fmt.Sprintf("driver factory %s already registered, ignoring", e.name)
}

var (
	driverMu         sync.RWMutex
	driverFactories  = make(map[string]factory)
	availableDrivers = map[string]factory{
		"pgx":   pgxV5Driver,
		"mysql": mysqlDriver,
	}
	errInvalidDriverName = fmt.Errorf(
		"invalid Driver name. Must be one of: %s",
		strings.Join(drivers(), ", "),
	)
)

func init() {
	if !testing.Testing() {
		registerAllDrivers()
	}
}

// Register registers a DB driver.
// Note: Register behaves similarly to database/sql.Register except that it doesn't
// panic on duplicate registrations, it just ignores them and continues.
func Register(name string, f factory) error {
	driverMu.Lock()
	defer driverMu.Unlock()

	if f == nil {
		panic(fmt.Sprintf("attempted to register driver %s with a nil factory", name))
	}

	if _, registered := driverFactories[name]; registered {
		return errFactoryAlreadyRegistered{name}
	}

	driverFactories[name] = f

	return nil
}

func registerAllDrivers() {
	for k, v := range availableDrivers {
		if err := Register(k, v); err != nil {
			panic(err)
		}
	}
}

func drivers() []string {
	driverMu.Lock()
	defer driverMu.Unlock()

	d := make([]string, 0, len(driverFactories))
	for k := range driverFactories {
		d = append(d, k)
	}

	slices.Sort(d)

	return d
}

// CreateDriver creates a Driver.
func CreateDriver(name string) (*Driver, error) {
	driverMu.Lock()
	defer driverMu.Unlock()

	driverFactory, ok := driverFactories[name]
	if !ok {
		return nil, errInvalidDriverName
	}

	return driverFactory(), nil
}

func mysqlDriver() *Driver {
	return &Driver{
		Driver:    &mysql.MySQLDriver{},
		Formatter: MysqlFormatter,
		AuthError: MySQLAuthError,
	}
}

func pgxV5Driver() *Driver {
	return &Driver{
		Driver:    &stdlib.Driver{},
		Formatter: PgFormatter,
		AuthError: PostgreSQLAuthError,
	}
}
