package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/musicmash/artisync/internal/db/models"
)

const (
	driver = "postgres"
)

type Conn struct {
	*models.Queries
	db *sql.DB
}

type Config struct {
	DSN                     string
	MaxOpenConnectionsCount int
	MaxIdleConnectionsCount int
	MaxConnectionIdleTime   time.Duration
	MaxConnectionLifetime   time.Duration
}

func Connect(conf Config) (*Conn, error) {
	db, err := sql.Open(driver, conf.DSN)
	if err != nil {
		return nil, fmt.Errorf("can't open connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping database: %w", err)
	}

	db.SetMaxOpenConns(conf.MaxOpenConnectionsCount)
	db.SetMaxIdleConns(conf.MaxIdleConnectionsCount)
	db.SetConnMaxIdleTime(conf.MaxConnectionIdleTime)
	db.SetConnMaxLifetime(conf.MaxConnectionLifetime)

	return &Conn{db: db, Queries: models.New(db)}, nil
}

func (conn *Conn) Ping() error {
	return conn.db.Ping()
}

func (conn *Conn) Close() error {
	return conn.db.Close()
}

func (conn *Conn) ExecTx(ctx context.Context, fn func(*models.Queries) error) error {
	tx, err := conn.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't begin tx: %w", err)
	}

	q := models.New(tx)
	if txErr := fn(q); txErr != nil {
		if err = tx.Rollback(); err != nil {
			//nolint:errorlint
			return fmt.Errorf("tx err: %v, rb err: %w", txErr, err)
		}

		return fmt.Errorf("pipeline failed: %w", err)
	}

	return tx.Commit()
}

func (conn *Conn) ApplyMigrations(filePath string) error {
	databaseInstance, err := postgres.WithInstance(conn.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can't create migrate postgres instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(filePath, driver, databaseInstance)
	if err != nil {
		return fmt.Errorf("can't create migrate file driver: %w", err)
	}

	err = m.Up()
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return fmt.Errorf("can't apply migrations: %w", err)
}
