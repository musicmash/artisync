package main

import (
	"bufio"
	"context"
	"os"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cron.Schedule(ctx, 5*time.Second, fetch)

	log.Info("artisync-daily is running...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
