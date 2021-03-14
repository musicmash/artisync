package xsync

import (
	"context"
	"database/sql"
	"errors"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/guard"
	"github.com/musicmash/artisync/internal/pipelines/syntask"
	"github.com/musicmash/artisync/internal/repository/sync"
	"golang.org/x/oauth2"
)

type Service struct {
	conn *db.Conn

	onceSync  syntask.Pipeline
	dailySync syntask.Pipeline
}

func NewService(conn *db.Conn, onceOAuthCredentials, dailyOAuthCredentials *oauth2.Config) *Service {
	service := Service{
		conn:      conn,
		onceSync:  syntask.New(onceOAuthCredentials, conn),
		dailySync: syntask.New(dailyOAuthCredentials, conn),
	}

	return &service
}

func convertPipelineResultToRepoTask(task *syntask.Task) *sync.Task {
	return &sync.Task{ID: task.ID}
}

func (s Service) DoOnceSync(ctx context.Context, userName, code string) (*sync.Task, error) {
	task, err := s.onceSync.Run(ctx, &syntask.PipelineOpts{
		UserName:          userName,
		SpotifyAuthCode:   code,
		ScheduleDailySync: false,
	})
	if err != nil {
		return nil, guard.NewInternalError(err)
	}

	return convertPipelineResultToRepoTask(task), nil
}

func (s Service) ConnectDailySync(ctx context.Context, userName, code string) (*sync.Task, error) {
	task, err := s.dailySync.Run(ctx, &syntask.PipelineOpts{
		UserName:          userName,
		SpotifyAuthCode:   code,
		ScheduleDailySync: true,
	})
	if err != nil {
		return nil, guard.NewInternalError(err)
	}

	return convertPipelineResultToRepoTask(task), nil
}

func (s Service) GetLatestSyncInfo(ctx context.Context, userName string) (*sync.LatestInfo, error) {
	info := sync.LatestInfo{}
	task, err := s.conn.GetLatestOneTimeSyncTask(ctx, userName)
	if err == nil {
		info.Latest = task.UpdatedAt.Format("2006-01-02T15:04:05")
		return &info, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		// no rows means that sync os disabled
		// return info with enabled: false (default bool value)
		return &info, nil
	}

	return nil, guard.NewInternalError(err)
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
