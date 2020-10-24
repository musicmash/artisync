package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/musicmash/artisync/internal/api/controllers/healthz"
	"github.com/musicmash/artisync/internal/api/controllers/spotify"
	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/pipelines/syntask"
)

func GetRouter(conn *db.Conn, pipeline syntask.Pipeline) chi.Router {
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

		r.Mount("/callbacks/spotify/artisync", spotify.New(pipeline).GetRouter())
	})

	return r
}