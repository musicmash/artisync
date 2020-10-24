package syntask

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/log"
)

type Task struct {
	ID      uuid.UUID `json:"id"`
	State   string    `json:"state"`
	Details *string   `json:"details"`
}

type Pipeline interface {
	GetOrCreateOneTimeTaskForUser(ctx context.Context, userName string, code string) (*Task, error)
	GetOrCreateDailyTaskForUser(ctx context.Context, userName string, code string) (*Task, error)
}

type TaskPipeline struct {
	// client *spotify.Client
	mgr *db.Conn
}

func New(mgr *db.Conn) Pipeline {
	return &TaskPipeline{mgr: mgr}
}

func (t *TaskPipeline) GetOrCreateOneTimeTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	// check if any tasks in-progress for user
	res, err := t.mgr.IsAnySyncTaskProcessingForUser(ctx, userName)
	if err == nil {
		// active task is found
		task := Task{ID: res.ID, State: string(res.State)}
		if res.Details.Valid {
			task.Details = &res.Details.String
		}

		return &task, nil
	}

	// task not found
	if errors.Is(err, sql.ErrNoRows) {
		// get user's refresh token

		task := models.ArtistOneTimeSyncTask{}
		err := t.mgr.ExecTx(ctx, func(db *models.Queries) error {
			task, err = db.CreateOneTimeSyncTask(ctx, userName)
			if err != nil {
				return guard.NewInternalError(fmt.Errorf("can't create sync task for %v: %w", userName, err))
			}

			// todo: replace when got access token
			params := models.CreateRefreshTokenParams{
				UserName:  userName,
				ExpiredAt: time.Now().Add(time.Hour),
				Value:     "fake-token",
			}
			if err = db.CreateRefreshToken(ctx, params); err != nil {
				return guard.NewInternalError(fmt.Errorf("can't save refresh token for %v: %w", userName, err))
			}

			return nil
		})
		if err != nil {
			return nil, guard.NewInternalError(fmt.Errorf("can't close tx: %w", err))
		}

		result := Task{ID: task.ID, State: string(task.State)}

		return &result, nil
	}

	// another error here
	return nil, guard.NewInternalError(fmt.Errorf("can't get cuurent tasks for %v: %w", userName, err))
}

func (t *TaskPipeline) GetOrCreateDailyTaskForUser(ctx context.Context, userName string, code string) (*Task, error) {
	// check if already tasks in-progress?
	// check if daily state is active
	log.Warn("will create only one time sync task, cause daily not implemented")

	return t.GetOrCreateOneTimeTaskForUser(ctx, userName, code)
}
