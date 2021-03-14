package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/artisync/internal/api/controllers/healthz"
	"github.com/musicmash/artisync/internal/api/controllers/sync"
	"github.com/musicmash/artisync/internal/api/controllers/tasks"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/repository"
)

func GetRouter(conn *db.Conn, repo *repository.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/_", func(r chi.Router) {
		r.Mount("/healthz", healthz.New(conn).GetRouter())
	})

	r.Route("/v1", func(r chi.Router) {
		// user logger inside /v1 route
		// to avoid logging of healthz requests
		r.Use(middleware.Logger)

		r.Mount("/artists/sync", sync.New(repo.Sync).GetRouter())
		r.Mount("/artists/sync/tasks", tasks.New(repo.Tasks).GetRouter())
	})

	return r
}
