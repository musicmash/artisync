package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/musicmash/artisync/internal/config"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/version"
)

var (
	configPath = flag.String("config", "", "abs path to conf file")
)

//nolint:funclen
func main() {
	flag.Parse()

	if *configPath == "" {
		_, _ = fmt.Fprintln(os.Stdout, "provide abs path to config via --config argument")
		return
	}

	conf, err := config.LoadFromFile(*configPath)
	if err != nil {
		exitIfError(err)
	}

	log.SetLevel(conf.Log.Level)
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)

	log.Info("connecting to db...")
	mgr, err := db.Connect(conf.DB.GetConnString())
	exitIfError(err)

	log.Info("connection to the db established")

	log.Info("applying migrations..")
	err = mgr.ApplyMigrations(conf.DB.MigrationsDir)
	if err != nil && err != migrate.ErrNoChange {
		exitIfError(fmt.Errorf("cant-t apply migrations: %v", err))
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
		exitIfError(fmt.Errorf("got error after tx: %s", err.Error()))
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(2)
}
