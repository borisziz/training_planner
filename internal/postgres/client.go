// Package postgres provides the application's PostgreSQL connection pool.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const defaultConnectTimeout = 5 * time.Second

const driverName = "pgx"

var ErrEmptyConnectionString = errors.New("postgres connection string is empty")

// Client is the SQLX database type used by repositories and transaction
// manager. PostgreSQL connections are created by the pgx driver.
type Client = *sqlx.DB

var _ trmsqlx.Tr = Client(nil)

// New creates an SQLX connection pool backed by pgx and verifies that
// PostgreSQL is reachable.
//
// connectionString accepts the PostgreSQL URL or keyword/value formats.
//
// The caller owns the returned client and must call Close during shutdown.
func New(ctx context.Context, connectionString string) (Client, error) {
	config, err := parseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	database := sqlx.NewDb(stdlib.OpenDB(*config), driverName)

	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return database, nil
}

func parseConfig(connectionString string) (*pgx.ConnConfig, error) {
	if strings.TrimSpace(connectionString) == "" {
		return nil, ErrEmptyConnectionString
	}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = defaultConnectTimeout
	}

	return config, nil
}
