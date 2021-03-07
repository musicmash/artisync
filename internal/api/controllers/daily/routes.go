package daily

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.GetDailySyncInfo)
	r.Delete("/", c.DisableDailySync)

	return r
}
