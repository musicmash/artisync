package daily

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
)

type Mgr struct {
	conn *db.Conn
}

func New(conn *db.Conn) *Mgr {
	return &Mgr{conn: conn}
}

type Info struct {
	Enabled  bool       `json:"enabled"`
	LastSync *time.Time `json:"last_sync"`
}

func (m *Mgr) GetDailyInfo(ctx context.Context, userName string) (*Info, error) {
	info := Info{}
	err := m.conn.ExecTx(ctx, func(db *models.Queries) error {
		_, err := db.GetUserDailySyncTask(ctx, userName)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("tried to get daily task for %s: %w", userName, err)
			}

			return nil
		}

		info.Enabled = true
		task, err := db.GetLatestOneTimeSyncTask(ctx, userName)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("tried to get latest one time sync task for %s: %w", userName, err)
			}

			return nil
		}

		info.LastSync = &task.CreatedAt
		return nil
	})
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return &info, nil
}

func (m *Mgr) DisableDailySync(ctx context.Context, userName string) error {
	err := m.conn.ExecTx(ctx, func(db *models.Queries) error {
		err := db.DisableDailySyncTask(ctx, userName)
		if err != nil {
			return fmt.Errorf("tried to disable daily sync task for %s: %w", userName, err)
		}

		return nil
	})
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	return nil
}
