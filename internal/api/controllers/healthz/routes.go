package healthz

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.Get)

	return r
}
