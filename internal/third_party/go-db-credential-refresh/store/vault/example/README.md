# Vault Example

This example is shows the Connector in operation. The only dependency required for this example is a running PostgreSQL instance. You can easily set one up by running:

```shell
docker run --name pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
```

Once the database is ready, you can run the example:

```shell
go run main.go
```

## How it works

Within `helpers.go` is a bunch of helper code that will create an in-memory Vault server, mount the database backend, and configure the postgres database and role.

Once the Vault server has been set up we instantiate the store:

```go
store, err := vault.NewStore(&vault.Config{
    Client:             client,
    CredentialLocation: vaultcredentials.NewAPIDatabaseCredentials(role, ""),
})
```

Note that the `CredentialLocation` is set to query the Vault API directly. The in-memory Vault server gives us a Vault client with the root token so we use that to authenticate.

If we were running in a Kubernetes cluster with the Vault k8s auth method, we could use the `vaultcredentials.AgentDatabaseCredentials` type to read the credentials directly from the agent-injected file. See the documentation on Vault's [Agent Sidecar Injector](https://www.vaultproject.io/docs/platform/k8s/injector) for more info.

There are a number of other Vault [auth methods](https://www.vaultproject.io/docs/auth) available although you'll need to write your own implementation.

Then we create the connector and set up the DB. For this example, we'll be using the [pgx](https://github.com/jackc/pgx) driver.

```go
c, err := driver.NewConnector(store, "pgx", &driver.Config{
    Host: host,
    Port: port,
    DB:   dbName,
    Opts: map[string]string{
        "sslmode": "disable",
    },
})

db := sql.OpenDB(c)
db.SetConnMaxLifetime(2 * time.Second)
db.SetMaxIdleConns(2)
db.SetMaxOpenConns(5)
```

With the `database/sql.DB` we can query the database, using the connector wrapped over the driver to manage retrieving dynamic credentials from Vault and using them to authenticate.

 See [`main.go`](main.go) for the rest of the example.
