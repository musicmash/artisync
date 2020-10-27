package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/musicmash/artisync/internal/cron"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/services/scheduletask"
	"github.com/musicmash/artisync/internal/version"
)

func main() {
	log.SetLevel("INFO")
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
