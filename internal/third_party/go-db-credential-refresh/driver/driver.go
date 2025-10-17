package driver

import (
	"database/sql/driver"
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/stdlib"
	v5 "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

// Driver carries information along with a database/sql/driver required for creating a Connector
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

//nolint:gochecknoglobals
var (
	driverMu         sync.RWMutex
	driverFactories  = make(map[string]factory)
	availableDrivers = map[string]factory{
		"pgxv4": pgxDriver,
		"pgx":   pgxV5Driver,
		"mysql": mysqlDriver,
		"pq":    pqDriver,
	}
	errInvalidDriverName = fmt.Errorf(
		"invalid Driver name. Must be one of: %s",
		strings.Join(drivers(), ", "),
	)
)

func init() { //nolint:gochecknoinits
	if !testing.Testing() {
		registerAllDrivers()
	}
}

// Register registers a DB driver
// Note: Register behaves similarly to database/sql.Register except that it doesn't
// panic on duplicate registrations, it just ignores them and continues.
// The reason we Register drivers separately from database/sql is because
//
//	a) most DB drivers already call database/sql.Register in an init() func
//	b) we need to carry a lot more information along with the driver to ensure our
//	   connector logic works correctly.
func Register(name string, f factory) error {
	driverMu.Lock()
	defer driverMu.Unlock()

	if f == nil {
		panic(fmt.Sprintf("attempted to register driver %s with a nil factory", name))
	}

	_, registered := driverFactories[name]
	if registered {
		return errFactoryAlreadyRegistered{name}
	}

	driverFactories[name] = f

	return nil
}

func registerAllDrivers() {
	for k, v := range availableDrivers {
		if err := Register(k, v); err != nil {
			// We should never, EVER hit this condition. If this happens it means something
			// has fundamentally broken in pgx, pq, or go-mysql.
			panic(err)
		}
	}
}

func drivers() []string {
	driverMu.Lock()
	defer driverMu.Unlock()

	drivers := make([]string, 0, len(driverFactories))
	for k := range driverFactories {
		drivers = append(drivers, k)
	}

	sort.Strings(drivers)

	return drivers
}

// CreateDriver creates a Driver.
func CreateDriver(name string) (*Driver, error) {
	driverMu.Lock()

	driverFactory, ok := driverFactories[name]
	if !ok {
		// Factory has not been registered.
		driverMu.Unlock()

		return nil, errInvalidDriverName
	}
	defer driverMu.Unlock()

	// Run the factory
	d := driverFactory()

	return d, nil
}

func mysqlDriver() *Driver {
	return &Driver{
		Driver:    &mysql.MySQLDriver{},
		Formatter: MysqlFormatter,
		AuthError: MySQLAuthError,
	}
}

func pgxDriver() *Driver {
	return &Driver{
		Driver:    &stdlib.Driver{},
		Formatter: PgFormatter,
		AuthError: PostgreSQLAuthError,
	}
}

func pgxV5Driver() *Driver {
	return &Driver{
		Driver:    &v5.Driver{},
		Formatter: PgFormatter,
		AuthError: PostgreSQLAuthError,
	}
}

func pqDriver() *Driver {
	return &Driver{
		Driver:    &pq.Driver{},
		Formatter: PgFormatter,
		AuthError: PostgreSQLAuthError,
	}
}
