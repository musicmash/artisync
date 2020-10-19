package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
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

	const pathToMigrations = "file:///etc/artisync/migrations"

	log.Info("applying migrations..")
	err = mgr.ApplyMigrations(pathToMigrations)
	if err != nil && err != migrate.ErrNoChange {
		log.Errorf("cant-t apply migrations: %v", err)
	}

	err = mgr.ExecTx(context.Background(), func(querier *models.Queries) error {
		art, err := querier.CreateArtist(context.Background(), models.CreateArtistParams{
			Name:   "rammstein",
			Poster: sql.NullString{},
		})
		if err != nil {
			return fmt.Errorf("can't create new artist: %w", err)
		}

		_, err = querier.CreateArtistAssociation(context.Background(), models.CreateArtistAssociationParams{
			ArtistID:  art.ID,
			StoreName: "spotify",
			StoreID:   "059c3940-a791-422d-8330-2954918c51e6",
		})
		if err != nil {
			return fmt.Errorf("can't associate artist: %w", err)
		}

		err = querier.CreateSubscription(context.Background(), models.CreateSubscriptionParams{
			UserName:  "objque",
			StoreName: "spotify",
			StoreID:   "059c3940-a791-422d-8330-2954918c51e6",
		})
		if err != nil {
			return fmt.Errorf("can't subscribe user: %w", err)
		}

		return nil
	})
	if err != nil {
		log.Errorf("got error after tx: %s", err.Error())
		os.Exit(2)
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
