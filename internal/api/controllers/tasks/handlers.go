package tasks

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/db"
)

type Controller struct {
	mgr *db.Conn
}

func New(mgr *db.Conn) *Controller {
	return &Controller{mgr: mgr}
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "task_id")
	if rawID == "" {
		httputils.WriteErrorWithCode(w, http.StatusNotFound, ErrTaskNotFound)
		return //nolint:nlreturn
	}

	taskID, err := uuid.Parse(rawID)
	if err != nil {
		httputils.WriteClientError(w, ErrInvalidUUID)
		return //nolint:nlreturn
	}

	state, err := c.mgr.GetOneTimeSyncTaskState(r.Context(), taskID)
	if err == nil {
		_ = httputils.WriteJSON(w, http.StatusOK, state)
		return //nolint:nlreturn
	}

	if errors.Is(err, sql.ErrNoRows) {
		httputils.WriteErrorWithCode(w, http.StatusNotFound, ErrTaskNotFound)
		return //nolint:nlreturn
	}

	httputils.WriteInternalError(w, err)
}
