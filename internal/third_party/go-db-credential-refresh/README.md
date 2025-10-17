# Go DB Credential Refresh

[![Godoc Reference](https://godoc.org/github.com/davepgreene/go-db-credential-refresh?status.svg)](https://pkg.go.dev/github.com/davepgreene/go-db-credential-refresh) 
[![Test](https://github.com/davepgreene/go-db-credential-refresh/workflows/Test/badge.svg)](https://github.com/davepgreene/go-db-credential-refresh/actions/workflows/test.yml) 
[![Lint](https://github.com/davepgreene/go-db-credential-refresh/workflows/Lint/badge.svg)](https://github.com/davepgreene/go-db-credential-refresh/actions/workflows/lint.yml) 
[![codecov](https://codecov.io/gh/davepgreene/go-db-credential-refresh/branch/master/graph/badge.svg)](https://codecov.io/gh/davepgreene/go-db-credential-refresh)

Go DB Credential Refresh is a driver to handle seamlessly reconnecting `database/sql` connections on credential 
rotation. This driver will work fine with static credentials but is designed for systems like 
[Hashicorp Vault](https://www.vaultproject.io/)'s 
[Database Secrets Engines](https://www.vaultproject.io/docs/secrets/databases) or 
[AWS RDS IAM Authentication](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html) 
where the credentials are retrieved from the identity manager before connecting.

Go DB Credential Refresh acts as a wrapper over existing DB drivers. It supports the following community DB 
drivers by default:

* [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)
* [pgx](https://github.com/jackc/pgx)
* [pq](https://github.com/lib/pq)

but users can register anything that implements 
[`database/sql/driver.Driver`](https://pkg.go.dev/database/sql/driver#Driver).

## Installation

```shell
go get -u github.com/davepgreene/go-db-credential-refresh
```

## Connector

The mechanism to interact with the driver is handled through a Connector which is a tight coupling between 
a `database/sql/driver.Driver`, a `Formatter`, and an `AuthError`. The latter two types handle formatting the 
components of a connection string for the specific DB implementation and an evaluation function that determines if 
an error coming from the `driver.Driver` is an authentication-related error.

## Formatters

`Formatters` assemble db- or driver-specific connection strings so the `Connector` can retry a connection with 
new credentials. This library ships with formatter implementations for MySQL and PostgreSQL both as a connection 
URI and a K/V connection string (see 
[the PostgreSQL docs](https://www.postgresql.org/docs/10/libpq-connect.html#LIBPQ-CONNSTRING) for more info) in 
the [`driver`](./driver) package.

## AuthErrors

An `AuthError` is an evaluative function which determines if an `error` represents a failed connection due to 
authentication. This tells the Connector to use its store to attempt to retrieve new credentials.`AuthError`s for 
MySQL and PostgreSQL are included in the `driver` package.

## Stores

A store is a mechanism to retrieve credentials. When you use the DB driver, you associate a `Store` with 
the `Connector`. Every time `Connector.Connect` is called, the store is queried for credentials. Stores must 
implement the `Store` interface (see [driver/store.go](driver/store.go)).

Go DB Credential Refresh currently ships with store implementations for Vault and RDS IAM Authentication. The 
Vault store includes both [Token Auth](https://www.vaultproject.io/docs/auth/token) and 
[Kubernetes Auth](https://www.vaultproject.io/docs/auth/kubernetes) authentication methods. See the 
[`vault`](./store/vault) package for the Vault implementation and [`awsrds`](./store/awsrds) package for RDS IAM
Authentication. Both included store implementations are available as independent modules.

## Examples

See the [examples directory](./examples) for sample usage and the Vault [example directory](./store/vault/example)
for how to use that module.
