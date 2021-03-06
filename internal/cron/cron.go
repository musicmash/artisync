package cron

import (
	"context"
	"time"

	"github.com/musicmash/artisync/internal/log"
)

type Task func(ctx context.Context) error

func Schedule(ctx context.Context, duration time.Duration, task Task) <-chan struct{} {
	log.Info("cron-job scheduled..")

	ticker := time.NewTicker(duration)

	done := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := task(ctx); err != nil {
					log.Error(err.Error())
				}
			case <-ctx.Done():
				log.Info("cron-job finished")
				ticker.Stop()
				done <- struct{}{}
				return
			}
		}
	}()

	return done
}
