package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/musicmash/artisync/internal/cron"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/version"
)

func fetch() error {
	log.Info("tick...")
	return nil
}

func main() {
	log.SetLevel("INFO")
	log.SetWriters(log.GetConsoleWriter())

	log.Debug(version.FullInfo)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	done := cron.Schedule(ctx, 5*time.Second, fetch)
	<-interrupt
	log.Info("got interrupt signal, shutdown..")
	cancel()

	<-done

	log.Info("daily-sync finished")
}
