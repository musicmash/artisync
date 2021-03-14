package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/musicmash/artisync/internal/api"
	"github.com/musicmash/artisync/internal/config"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/repository"
	"github.com/musicmash/artisync/internal/services/spotify/auth"
	"github.com/musicmash/artisync/internal/services/tasks"
	"github.com/musicmash/artisync/internal/services/xsync"
	"github.com/musicmash/artisync/internal/version"
)

//nolint:funclen
func main() {
	_ = flag.Bool("version", false, "show build info and exit")
	if versionRequired() {
		_, _ = fmt.Fprintln(os.Stdout, version.FullInfo)
		os.Exit(0)
	}

	_ = flag.Bool("help", false, "show this message and exit")
	if helpRequired() {
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := flag.String("config", "", "abs path to conf file")
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
	mgr, err := db.Connect(db.Config{
		DSN:                     conf.DB.GetConnString(),
		MaxOpenConnectionsCount: conf.DB.MaxOpenConnections,
		MaxIdleConnectionsCount: conf.DB.MaxIdleConnections,
		MaxConnectionIdleTime:   conf.DB.MaxConnectionIdleTime,
		MaxConnectionLifetime:   conf.DB.MaxConnectionLifeTime,
	})
	exitIfError(err)

	log.Info("connection to the db established")

	if conf.DB.AutoMigrate {
		log.Info("applying migrations..")
		exitIfError(mgr.ApplyMigrations(conf.DB.MigrationsDir))
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), conf.HTTP.WriteTimeout)
	defer cancel()

	exitIfError(auth.ValidateAuthConf(&conf.Spotify))

	repo := repository.Repository{
		Sync: xsync.NewService(
			mgr,
			conf.Spotify.GetOnceSyncOAuthConfig(),
			conf.Spotify.GetDailySyncOAuthConfig(),
		),
		Tasks: tasks.New(mgr),
	}
	router := api.GetRouter(mgr, &repo)
	server := api.New(router, conf.HTTP)

	go gracefulShutdown(ctx, server, quit, done)

	log.Infof("server is ready to handle requests at: %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		exitIfError(fmt.Errorf("could not listen on %v: %v", server.Addr, err))
	}

	<-done
	_ = mgr.Close()
	log.Info("artisync-api stopped")
}

func gracefulShutdown(ctx context.Context, server *api.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Info("server is shutting down...")

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

func isArgProvided(argName string) bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, argName) {
			return true
		}
	}
	return false
}

func helpRequired() bool {
	return isArgProvided("-help")
}

func versionRequired() bool {
	return isArgProvided("-version")
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(2)
}
