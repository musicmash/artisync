package sync

import (
	"fmt"
	"net/http"

	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/repository/sync"
)

type Controller struct {
	repo sync.Repository
}

func New(repo sync.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) processSpotifyCallback(w http.ResponseWriter, r *http.Request, isDailyCallback bool) {
	var (
		values   = r.URL.Query()
		userName = httputils.GetUserName(r)
		code     = values.Get("code")

		task *sync.Task
		err  error
	)

	if err := values.Get("error"); err != "" {
		if err == "access_denied" {
			log.Infof("user '%s' denied access to Spotify", userName)
		} else {
			log.Errorf("got error query '%s' for user '%s'", err, userName)
		}

		http.Redirect(w, r, "/subscriptions", http.StatusMovedPermanently)
		return
	}

	if isDailyCallback {
		task, err = c.repo.ConnectDailySync(r.Context(), userName, code)
	} else {
		task, err = c.repo.DoOnceSync(r.Context(), userName, code)
	}
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	url := fmt.Sprintf("/subscriptions?task_id=%s", task.ID.String())
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (c *Controller) DoOnceSync(w http.ResponseWriter, r *http.Request) {
	const isDailyCallback = false
	c.processSpotifyCallback(w, r, isDailyCallback)
}

func (c *Controller) ConnectDailySync(w http.ResponseWriter, r *http.Request) {
	const isDailyCallback = true
	c.processSpotifyCallback(w, r, isDailyCallback)
}

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
