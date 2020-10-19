package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/version"
)

func main() {
	log.SetLevel("DEBUG")
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	log.Info("connecting to db...")
	mgr, err := db.Connect(dsn)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	log.Info("connection to the db established")

	err = mgr.ExecTx(context.Background(), func(querier *models.Queries) error {
		art, err := querier.CreateArtist(context.Background(), models.CreateArtistParams{
			Name:   "rammstein",
			Poster: sql.NullString{},
		})
		if err != nil {
			return err
		}

		log.Infof("id: %d, created_at: %v", art.ID, art.CreatedAt.String())
		return nil
	})
	if err != nil {
		log.Errorf("got error after tx: %s", err.Error())
		os.Exit(2)
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
