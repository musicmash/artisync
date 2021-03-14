package repository

import (
	"github.com/musicmash/artisync/internal/repository/sync"
	"github.com/musicmash/artisync/internal/repository/tasks"
)

type Repository struct {
	Sync  sync.Repository
	Tasks tasks.Repository
}
