package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
	log.Println("connecting to db...")
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	log.Println("connection to the db established!")
	return &Conn{db: db, Queries: models.New(db)}, nil
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
