package spotify

import (
	"fmt"
	"net/http"

	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/log"
	"github.com/musicmash/artisync/internal/services/syntask"
)

type Controller struct {
	mgr *syntask.Mgr
}

func New(mgr *syntask.Mgr) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) ProcessCallback(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	if err := values.Get("error"); err != "" {
		if err != "access_denied" {
			log.Errorf("got '%v' error query when try to sync artists", err)
		}

		url := fmt.Sprintf("/onboarding/artist-sync?error=%v", err)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}

	if err := validateStateAndCode(values); err != nil {
		httputils.WriteClientError(w, err)
		return
	}

	var (
		userName = httputils.GetUserName(r)
		code     = values.Get("code")

		task *syntask.Task
		err  error
	)

	isDaily := values.Get("state") == stateBackgroundSyncAllowed
	if isDaily {
		task, err = c.mgr.GetOrCreateDailySyncTaskForUser(r.Context(), userName, code)
	} else {
		task, err = c.mgr.GetOrCreateOneTimeSyncTaskForUser(r.Context(), userName, code)
	}
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	url := fmt.Sprintf("/subscriptions?task_id=%v", task.ID.String())
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
