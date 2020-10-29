package sync

import (
	"context"

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

func (mgr *Mgr) WaitUntilAllWorkesFinish() {
	var workerID uint8
	for workerID = 1; workerID <= mgr.conf.WorkersCount; workerID++ {
		<-mgr.done
	}

	log.Infof("successfully stop %d sync workers", mgr.conf.WorkersCount)
}

func (mgr *Mgr) Schedule(ctx context.Context) error {
	// get all tasks with state=created
	// set state=scheduled
	// send their to workers
	log.Info("looking new tasks with state=created")
	mgr.tasks <- &models.GetNextOneTimeSyncTasksRow{}
	return nil
}
