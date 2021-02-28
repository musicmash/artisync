package daily

import (
	"net/http"

	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/services/daily"
)

type Controller struct {
	mgr *daily.Mgr
}

func New(mgr *daily.Mgr) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) GetDailySyncInfo(w http.ResponseWriter, r *http.Request) {
	userName := httputils.GetUserName(r)
	info, err := c.mgr.GetDailyInfo(r.Context(), userName)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, &info)
}

func (c *Controller) DisableDailySync(w http.ResponseWriter, r *http.Request) {
	userName := httputils.GetUserName(r)
	err := c.mgr.DisableDailySync(r.Context(), userName)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
