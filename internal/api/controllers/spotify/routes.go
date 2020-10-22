package spotify

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.ArtistsSyncCallback)

	return r
}
