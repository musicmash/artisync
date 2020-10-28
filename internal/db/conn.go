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
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
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
		return err
	}

	q := models.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

func (conn *Conn) ApplyMigrations(filePath string) error {
	driver, err := postgres.WithInstance(conn.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can't create migrate postgres instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres",
		driver,
	)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("can't create migrate file driver: %w", err)
	}

	return m.Up()
}
