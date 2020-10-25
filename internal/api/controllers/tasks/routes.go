package tasks

import "github.com/go-chi/chi"

func (c *Controller) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/{task_id}", c.GetTask)

	return r
}
