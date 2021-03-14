package sync

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/once/connect", c.DoOnceSync)
	r.Get("/daily/connect", c.ConnectDailySync)

	r.Get("/daily", c.GetDailySyncInfo)
	r.Delete("/daily", c.DisableDailySync)

	return r
}
