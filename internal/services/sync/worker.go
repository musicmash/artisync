package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/pipelines/sync"
)

func updateTaskState(ctx context.Context, db *db.Conn, taskID uuid.UUID, state models.TaskState) error {
	err := db.UpdateOneTimeSyncTaskState(ctx, models.UpdateOneTimeSyncTaskStateParams{
		ID:      taskID,
		State:   state,
		Details: []byte("{}"),
	})
	if err != nil {
		return fmt.Errorf("can't update state for task: %s: %w", taskID.String(), err)
	}

	return nil
}

func (mgr *Mgr) processTask(ctx context.Context, task *models.GetNextOneTimeSyncTasksRow) error {
	// set state as in-progress
	if err := updateTaskState(ctx, mgr.conn, task.ID, models.TaskStateInProgress); err != nil {
		return fmt.Errorf("can't update state for task with id %s: %w", task.ID.String(), err)
	}

	// run pipeline that syncs subscriptions
	err := mgr.pipeline.Run(ctx, &sync.PipelineOpts{
		UserName:     task.UserName,
		RefreshToken: task.RefreshToken,
	})
	// update state set err if pipeline failed
	if err != nil {
		if err := updateTaskState(ctx, mgr.conn, task.ID, models.TaskStateError); err != nil {
			return fmt.Errorf("can't update state for failed task with id %s: %w", task.ID.String(), err)
		}

		return fmt.Errorf("pipeline failed for task with id %s: %w", task.ID.String(), err)
	}

	// update state set done if pipeline finished without errors
	if err := updateTaskState(ctx, mgr.conn, task.ID, models.TaskStateDone); err != nil {
		return fmt.Errorf("can't update state for finished task with id %s", task.ID.String())
	}

	return nil
}

func (mgr *Mgr) runWorker(ctx context.Context, workerID uint8) {
	log.Infof("sync worker #%d is running...", workerID)
	for {
		select {
		case task := <-mgr.tasks:
			start := time.Now().UTC()
			log.Infof("#%d worker got task, processing it...", workerID)
			if err := mgr.processTask(ctx, task); err != nil {
				log.Error(err.Error())
			}
			elapsed := time.Now().UTC().Sub(start)
			log.Infof("#%d worker finish processing task, elapsed %s", workerID, elapsed.String())
		case <-ctx.Done():
			log.Infof("sync worker #%d is finished", workerID)
			mgr.done <- struct{}{}
			return
		}
	}
}
