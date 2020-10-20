package db

import (
	"context"
	"database/sql"
	"fmt"

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

func Connect(dsn string) (*Conn, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

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
	if err != nil {
		return fmt.Errorf("can't create migrate file driver: %w", err)
	}

	return m.Up()
}
