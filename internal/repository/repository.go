package repository

import "github.com/musicmash/artisync/internal/repository/sync"

type Repository struct {
	Sync sync.Repository
}
