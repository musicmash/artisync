package sync

import (
	"net/http"

	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/repository/sync"
)

type Controller struct {
	repo sync.Repository
}

func New(repo sync.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) DoOnceSync(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) ConnectDailySync(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) GetLatestSyncInfo(w http.ResponseWriter, r *http.Request) {
	userName := httputils.GetUserName(r)

	info, err := c.repo.GetLatestSyncInfo(r.Context(), userName)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &info)
}

func (c *Controller) GetDailySyncInfo(w http.ResponseWriter, r *http.Request) {
	userName := httputils.GetUserName(r)

	info, err := c.repo.GetDailySyncInfo(r.Context(), userName)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &info)
}

func (c *Controller) DisableDailySync(w http.ResponseWriter, r *http.Request) {
	userName := httputils.GetUserName(r)

	err := c.repo.DisableDailySync(r.Context(), userName)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
