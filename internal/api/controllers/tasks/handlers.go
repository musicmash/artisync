package tasks

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/services/syntask"
)

type Controller struct {
	mgr *syntask.Mgr
}

func New(mgr *syntask.Mgr) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "task_id")
	if rawID == "" {
		httputils.WriteErrorWithCode(w, http.StatusNotFound, ErrTaskNotFound)
		return
	}

	taskID, err := uuid.Parse(rawID)
	if err != nil {
		httputils.WriteClientError(w, ErrInvalidUUID)
		return
	}

	state, err := c.mgr.GetSyncTaskState(r.Context(), taskID)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, state)
}
