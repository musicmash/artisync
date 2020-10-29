package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/musicmash/artisync/internal/config"
	"github.com/musicmash/artisync/internal/cron"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/log"
	pipeline "github.com/musicmash/artisync/internal/pipelines/sync"
	"github.com/musicmash/artisync/internal/services/sync"
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
		err = mgr.ApplyMigrations(conf.DB.MigrationsDir)
		if err != nil {
			exitIfError(fmt.Errorf("cant-t apply migrations: %v", err))
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	pipe := pipeline.New(mgr)
	task := sync.New(mgr, pipe, sync.WorkerConfig{
		// TODO (m.kalinin): extract values into config
		WorkersCount: 5,
		TasksCount:   30,
	})
	task.RunWorkers(ctx)
	done := cron.Schedule(ctx, 10*time.Second, task.Schedule)

	<-interrupt
	log.Info("got interrupt signal, shutdown..")
	cancel()

	<-done
	task.WaitWorkers()

	log.Info("artisync-sync finished")
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
