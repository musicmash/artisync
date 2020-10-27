package cron

import (
	"context"
	"time"

	"github.com/musicmash/artisync/internal/log"
)

type Task func() error

func Schedule(ctx context.Context, duration time.Duration, task Task) {
	log.Info("cron-job scheduled..")

	ticker := time.NewTicker(duration)
	defer func() {
		log.Info("cron-job cancelled")
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if err := task(); err != nil {
				log.Error(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}
