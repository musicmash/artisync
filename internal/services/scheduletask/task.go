package scheduletask

import (
	"context"
	"fmt"
	"time"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/log"
)

type Task struct {
	conn *db.Conn
}

func New(conn *db.Conn) *Task {
	return &Task{conn: conn}
}

func (t *Task) Schedule(ctx context.Context) error {
	return t.conn.ExecTx(ctx, func(db *models.Queries) error {
		// get lock to avoid race condition between instances
		if err := db.GetDailyLock(ctx); err != nil {
			return fmt.Errorf("can't get lock: %w", err)
		}

		// schedule tasks for users that weren't updated and whose token is still alive
		today := time.Now().UTC().Truncate(24 * time.Hour)
		count, err := db.ScheduleDailyTasks(ctx, models.ScheduleDailyTasksParams{
			Today:     today,
			Yesterday: today.Add(-24 * time.Hour),
		})
		if err != nil {
			return fmt.Errorf("can't schedule daily tasks: %w", err)
		}

		log.Infof("scheduled %v daily tasks", count)

		if count == 0 {
			return nil
		}

		// reset updated_at
		count, err = db.ResetDailyTasks(ctx, time.Now().UTC().Truncate(24*time.Hour))
		if err != nil {
			return fmt.Errorf("can't update daily tasks: %w", err)
		}

		log.Infof("successfully update %v daily tasks", count)
		return nil
	})
}
