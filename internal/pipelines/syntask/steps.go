package syntask

import (
	"context"
	"fmt"

	"github.com/musicmash/artisync/internal/db/models"
	"github.com/musicmash/artisync/internal/guard"
)

func GetRefreshTokenStep(ctx context.Context, opts *PipelineOpts, data *PipelineData) error {
	token, err := opts.auth.Exchange(ctx, opts.SpotifyAuthCode)
	if err != nil {
		return fmt.Errorf("can't get access_token: %w", err)
	}

	data.token = token
	return nil
}

func ScheduleSyncTaskStep(ctx context.Context, opts *PipelineOpts, data *PipelineData) error {
	err := data.conn.ExecTx(ctx, func(db *models.Queries) error {
		task, err := db.CreateOneTimeSyncTask(ctx, opts.UserName)
		if err != nil {
			return fmt.Errorf("can't create one-time sync task for %v: %w", opts.UserName, err)
		}

		if opts.ScheduleDailySync {
			if err := db.CreateDailySyncTask(ctx, opts.UserName); err != nil {
				return fmt.Errorf("can't create daily sync task for %v: %w", opts.UserName, err)
			}
		}

		params := models.CreateRefreshTokenParams{
			UserName:  opts.UserName,
			ExpiredAt: data.token.Expiry,
			Value:     data.token.RefreshToken,
		}
		if err = db.CreateRefreshToken(ctx, params); err != nil {
			return fmt.Errorf("can't save refresh token for %v: %w", opts.UserName, err)
		}

		data.task = &Task{ID: task.ID, State: string(task.State)}
		return nil
	})
	if err != nil {
		return guard.NewInternalError(fmt.Errorf("can't close tx: %w", err))
	}

	return nil
}
