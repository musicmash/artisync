package tasks

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/musicmash/artisync/internal/api/httputils"
	"github.com/musicmash/artisync/internal/repository/tasks"
)

type Controller struct {
	repo tasks.Repository
}

func New(repo tasks.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) GetTask(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "task_id")
	if rawID == "" {
		httputils.WriteErrorWithCode(w, http.StatusBadRequest, ErrTaskNotFound)
		return
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		httputils.WriteClientError(w, ErrInvalidUUID)
		return
	}

	task, err := c.repo.GetTask(r.Context(), id)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	_ = httputils.WriteJSON(w, http.StatusOK, task)
}
