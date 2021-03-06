// Code generated by sqlc. DO NOT EDIT.

package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskState string

const (
	TaskStateCreated    TaskState = "created"
	TaskStateScheduled  TaskState = "scheduled"
	TaskStateInProgress TaskState = "in-progress"
	TaskStateDone       TaskState = "done"
	TaskStateError      TaskState = "error"
)

func (e *TaskState) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TaskState(s)
	case string:
		*e = TaskState(s)
	default:
		return fmt.Errorf("unsupported scan type for TaskState: %T", src)
	}
	return nil
}

type Artist struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	Name      string         `json:"name"`
	Poster    sql.NullString `json:"poster"`
}

type ArtistAssociation struct {
	ID        int32  `json:"id"`
	ArtistID  int64  `json:"artist_id"`
	StoreName string `json:"store_name"`
	StoreID   string `json:"store_id"`
}

type ArtistDailySyncTask struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
}

type ArtistOneTimeSyncTask struct {
	ID        uuid.UUID       `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	UserName  string          `json:"user_name"`
	State     TaskState       `json:"state"`
	Details   json.RawMessage `json:"details"`
}

type ArtistSyncRefreshToken struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	UserName  string    `json:"user_name"`
	Value     string    `json:"value"`
}

type Subscription struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name"`
	ArtistID  int64     `json:"artist_id"`
}
