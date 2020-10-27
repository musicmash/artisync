package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/musicmash/artisync/internal/config"
	"github.com/musicmash/artisync/internal/cron"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/services/scheduletask"
	"github.com/musicmash/artisync/internal/version"
)

func main() {
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

	task := scheduletask.New()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	done := cron.Schedule(ctx, 5*time.Second, task.Schedule)
	<-interrupt
	log.Info("got interrupt signal, shutdown..")
	cancel()

	<-done

	log.Info("daily-sync finished")
}

func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(2)
}
