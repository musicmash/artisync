package sync

import (
	"context"
	"fmt"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/pipelines/sync"
)

type Mgr struct {
	conn     *db.Conn
	pipeline sync.Pipeline
	conf     WorkerConfig
	tasks    chan *models.GetNextOneTimeSyncTasksRow
	done     chan struct{}
}

type WorkerConfig struct {
	WorkersCount uint8
	TasksCount   uint16
}

func New(conn *db.Conn, pipeline sync.Pipeline, conf WorkerConfig) *Mgr {
	return &Mgr{
		conn:     conn,
		pipeline: pipeline,
		conf:     conf,
		tasks:    make(chan *models.GetNextOneTimeSyncTasksRow, conf.TasksCount),
		done:     make(chan struct{}, conf.WorkersCount),
	}
}

func (mgr *Mgr) RunWorkers(ctx context.Context) {
	var workerID uint8
	for workerID = 1; workerID <= mgr.conf.WorkersCount; workerID++ {
		go mgr.runWorker(ctx, workerID)
	}

	log.Infof("successfully spawn %d sync workers", mgr.conf.WorkersCount)
}

func (mgr *Mgr) WaitWorkers() {
	var workerID uint8
	for workerID = 1; workerID <= mgr.conf.WorkersCount; workerID++ {
		<-mgr.done
	}

	log.Infof("successfully stop %d sync workers", mgr.conf.WorkersCount)
}

func (mgr *Mgr) Schedule(ctx context.Context) error {
	err := mgr.conn.ExecTx(ctx, func(db *models.Queries) error {
		// get tasks with state=created
		const tasksLimit = 3
		tasks, err := db.GetNextOneTimeSyncTasks(ctx, tasksLimit)
		if err != nil {
			return fmt.Errorf("can't get tasks with state=created for update: %w", err)
		}

		if len(tasks) == 0 {
			log.Warn("nothing to sync :(")
			return nil
		}

		// set state as scheduled
		for i := range tasks {
			err := db.UpdateOneTimeSyncTaskState(ctx, models.UpdateOneTimeSyncTaskStateParams{
				ID:      tasks[i].ID,
				State:   models.TaskStateScheduled,
				Details: []byte("{}"),
			})
			if err != nil {
				return fmt.Errorf("can't update state for task: %s: %w", tasks[i].ID.String(), err)
			}

			// send tasks to workers
			mgr.tasks <- &tasks[i]
		}

		log.Infof("scheduled %d tasks", len(tasks))
		return nil
	})
	if err != nil {
		return fmt.Errorf("can't spawn tasks with state=created: %w", err)
	}

	return nil
}
