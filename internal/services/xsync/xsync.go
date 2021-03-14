package xsync

import (
	"context"
	"database/sql"
	"errors"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/repository/sync"
)

type Service struct {
	conn *db.Conn
}

func NewService(conn *db.Conn) *Service {
	return &Service{conn: conn}
}

func (s Service) DoOnceSync(ctx context.Context) error {
	panic("implement me")
}

func (s Service) ConnectDailySync(ctx context.Context) error {
	panic("implement me")
}

func (s Service) GetDailySyncInfo(ctx context.Context, userName string) (*sync.DailyInfo, error) {
	info := sync.DailyInfo{}
	_, err := s.conn.GetUserDailySyncTask(ctx, userName)
	if err == nil {
		info.Enabled = true
		return &info, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		// no rows means that sync os disabled
		// return info with enabled: false (default bool value)
		return &info, nil
	}

	return nil, guard.NewInternalError(err)
}

func (s Service) DisableDailySync(ctx context.Context, userName string) error {
	err := s.conn.DisableDailySyncTask(ctx, userName)
	if err != nil {
		return guard.NewInternalError(err)
	}

	return nil
}
