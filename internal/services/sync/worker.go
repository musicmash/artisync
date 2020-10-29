package sync

import (
	"context"
	"time"

	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
)

func (mgr *Mgr) processTask(task *models.GetNextOneTimeSyncTasksRow) {
	// set state as in-progress
	// run pipeline
	// if err != nil
	// set state as err
	// return
	// set state as done
	time.Sleep(5 * time.Second)
}

func (mgr *Mgr) runWorker(ctx context.Context, workerID uint8) {
	log.Infof("sync worker #%d is running...", workerID)
	for {
		select {
		case task := <-mgr.tasks:
			start := time.Now().UTC()
			log.Infof("#%d worker got task, processing it...", workerID)
			mgr.processTask(task)
			elapsed := time.Now().UTC().Sub(start)
			log.Infof("#%d worker finish processing task, elapsed %s", workerID, elapsed.String())
		case <-ctx.Done():
			log.Infof("sync worker #%d is finished", workerID)
			mgr.done <- struct{}{}
			return
		}
	}
}
