package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"port-location/internal/portdomain"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

/*
DBClient implemented as fast solution for integration tests.
It would be better to control spawning and migrating db from test code.
Also it would be nice to add a mechanism for populating db with testdata from fixtures (CSV files for example)
*/

type DBClient struct {
	DB   *sqlx.DB
	conf portdomain.DB
}

func NewTestDBClient(t *testing.T) DBClient {
	testConf := portdomain.DB{
		Host:          "localhost",
		Port:          "5435",
		Name:          "postgres",
		Username:      "postgres",
		Password:      "example",
		MigrationsDir: os.Getenv("MIGRATIONS_PATH"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	psqlConf := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testConf.Host, testConf.Port, testConf.Username, testConf.Password, testConf.Name,
	)

	dbs, err := sqlx.ConnectContext(ctx, "postgres", psqlConf)
	require.NoError(t, err)

	d := DBClient{
		conf: testConf,
		DB:   dbs,
	}

	d.mustMigrate(t, true)

	return d
}

func (td *DBClient) mustMigrate(t *testing.T, migrateUp bool) {
	driver, err := postgres.WithInstance(td.DB.DB, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(td.conf.MigrationsDir, td.conf.Name, driver)
	require.NoError(t, err)

	var migrateErr error

	if migrateUp {
		migrateErr = m.Up()
	} else {
		migrateErr = m.Down()
	}

	if migrateErr != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
}

func (td *DBClient) Close(t *testing.T) {
	td.mustMigrate(t, false)
	if err := td.DB.Close(); err != nil {
		require.NoError(t, err)
	}
}
