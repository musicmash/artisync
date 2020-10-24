package spotify

import (
	"fmt"
	"net/http"

	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/pipelines/syntask"
)

type Controller struct {
	pipeline syntask.Pipeline
}

func New(pipeline syntask.Pipeline) *Controller {
	return &Controller{pipeline: pipeline}
}

func (c *Controller) OneTimeSyncCallback(w http.ResponseWriter, r *http.Request) {
	c.processCallback(false, w, r)
}

func (c *Controller) DailySyncCallback(w http.ResponseWriter, r *http.Request) {
	c.processCallback(true, w, r)
}

func (c *Controller) processCallback(isDaily bool, w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if err := validateQuery(r.URL.Query()); err != nil {
		httputils.WriteClientError(w, err)
		return //nolint:nlreturn
	}

	var (
		userName = httputils.GetUserName(r)
		code     = values.Get("code")

		task *syntask.Task
		err  error
	)

	if isDaily {
		task, err = c.pipeline.GetOrCreateDailyTaskForUser(r.Context(), userName, code)
	} else {
		task, err = c.pipeline.GetOrCreateSingleTaskForUser(r.Context(), userName, code)
	}
	if err != nil {
		httputils.WriteGuardError(w, err)
		return //nolint:nlreturn
	}

	w.WriteHeader(http.StatusCreated)
	url := fmt.Sprintf("/onboarding/artist-sync?task_id=%v", task.ID.String())
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
