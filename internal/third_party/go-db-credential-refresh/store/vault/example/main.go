package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/davepgreene/go-db-credential-refresh/examples/db"

	"github.com/davepgreene/go-db-credential-refresh/store/vault"
	vaultcredentials "github.com/davepgreene/go-db-credential-refresh/store/vault/credentials"
)

const (
	role     = "role"
	host     = "localhost"
	username = "postgres"
	password = "postgres"
	dbName   = "postgres"
	port     = 5432
)

func main() {
	err := Run()
	if err == nil {
		return
	}

	if err == context.Canceled {
		return
	}

	log.Fatal(err)
}

func Run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	infraCtx := context.Background()

	dbContainer, err := setupDb(infraCtx)
	if err != nil {
		return err
	}

	defer func() {
		if cleanupErr := dbContainer.Terminate(infraCtx); cleanupErr != nil {
			log.Printf("error cleaning up postgres container: %v", cleanupErr)
			err = cleanupErr
		}
	}()

	driverConfig, err := generateDriverConfig(infraCtx, dbContainer)
	if err != nil {
		return err
	}

	// Set up Vault, DB backend, and Postgres configuration
	client, cleanup, err := setupVault(infraCtx, "host.docker.internal", driverConfig.Port)
	if err != nil {
		if cleanup != nil {
			if err := cleanup(); err != nil {
				return err
			}
		}

		return err
	}
	defer func() {
		if cleanupErr := cleanup(); cleanupErr != nil {
			log.Printf("error cleaning up vault container: %v", cleanupErr)
			err = cleanupErr
		}
	}()

	// Create the store
	store, err := vault.NewStore(&vault.Config{
		Client:             client,
		CredentialLocation: vaultcredentials.NewAPIDatabaseCredentials(role, ""),
	})
	if err != nil {
		return err
	}

	// Create the connector which implements database/sql/driver.Connector
	c, err := driver.NewConnector(store, "pgx", driverConfig)
	if err != nil {
		return err
	}

	// Use the built in database/sql package to work with the connector
	database := sql.OpenDB(c)

	// In order to demonstrate the creation and revocation of roles we need to set the
	// connection lifetime very short. In a production environment, Vault role TTLs and
	// connection lifetime should be tuned based on database performance requirements.
	database.SetConnMaxLifetime(2 * time.Second)
	database.SetMaxIdleConns(2)
	database.SetMaxOpenConns(5)

	appSignal := make(chan os.Signal, 1)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		cancel()
	}()

	for {
		// First ping the DB to open a connection
		if err := db.Ping(ctx, database); err != nil {
			log.Println(err)

			break
		}

		// Sleep long enough that the creds should expire
		time.Sleep(3 * time.Second)

		// Now get users
		users, err := db.QueryUsers(ctx, database, map[string]bool{
			username: false,
		})
		if err != nil {
			fmt.Println(err)

			break
		}

		fmt.Printf("Retrieving users from database: %v\n", users)
	}

	// err = TearDownRoles(ctx, client)

	return err
}
