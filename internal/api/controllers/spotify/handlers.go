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

func (c *Controller) ArtistsSyncCallback(w http.ResponseWriter, r *http.Request) {
	// check is code provided
	// check state
	//

	const (
		userName = "objque"
		state    = "once-sync"
	)

	task, err := c.pipeline.GetOrCreateSingleTaskForUser(r.Context(), userName, state)
	if err != nil {
		httputils.WriteGuardError(w, err)
		return
	}

	fmt.Fprintf(w, task.ID.String())
}
