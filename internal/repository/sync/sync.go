package sync

import (
	"context"

	"github.com/google/uuid"
)

type LatestInfo struct {
	Latest string `json:"latest"`
}

type DailyInfo struct {
	Enabled bool `json:"enabled"`
}

type Task struct {
	ID uuid.UUID `json:"id"`
}

type Repository interface {
	// Sync methods here
	DoOnceSync(ctx context.Context, userName, code string) (*Task, error)
	ConnectDailySync(ctx context.Context, userName, code string) (*Task, error)

	// Daily methods here
	GetDailySyncInfo(ctx context.Context, userName string) (*DailyInfo, error)
	DisableDailySync(ctx context.Context, userName string) error

	// Other methods here
	GetLatestSyncInfo(ctx context.Context, userName string) (*LatestInfo, error)
}
