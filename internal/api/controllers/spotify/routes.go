package spotify

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/one-time", c.OneTimeSyncCallback)
	r.Get("/daily", c.DailySyncCallback)

	return r
}
