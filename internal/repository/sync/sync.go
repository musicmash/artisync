package sync

import "context"

type DailyInfo struct {
	Enabled bool `json:"enabled"`
}

type Repository interface {
	// Sync methods here
	DoOnceSync(ctx context.Context) error
	ConnectDailySync(ctx context.Context) error

	// Daily methods here
	GetDailySyncInfo(ctx context.Context, userName string) (*DailyInfo, error)
	DisableDailySync(ctx context.Context, userName string) error
}
