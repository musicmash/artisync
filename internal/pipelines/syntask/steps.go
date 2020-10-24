package syntask

import (
	"context"
	"fmt"

	"github.com/musicmash/artisync/internal/db/models"
)

func GetRefreshTokenStep(ctx context.Context, opts *PipelineOpts, data *PipelineData) error {
	token, err := data.auth.Exchange(ctx, opts.SpotifyAuthCode)
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
		} else {
			// by default token.Expiry shows us time
			// when access-token will expire.
			// refresh-token keep alive until user manually
			// disconnect our app in the settings.
			//
			// so, if user wanna periodically sync his artists
			// we should override/prolong token.Expiry time
			// e.g for 3 months
			data.token.Expiry = data.token.Expiry.AddDate(0, 3, 0)
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
		return fmt.Errorf("can't close tx: %w", err)
	}

	return nil
}
