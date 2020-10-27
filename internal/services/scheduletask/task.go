package scheduletask

import (
	"github.com/musicmash/artisync/internal/log"
)

type Task struct{}

func New() *Task {
	return &Task{}
}

func (t *Task) Schedule() error {
	// get lock on the db
	// get users that weren't updated and whose token is still alive
	// schedule tasks for them
	// commit tx
	log.Info("tick...")
	return nil
}
