package spotify

import (
	"fmt"
	"net/http"

	"github.com/musicmash/artisync/internal/pipelines/syntask"
)

type Controller struct {
	pipeline syntask.Pipeline
}

func New(pipeline syntask.Pipeline) *Controller {
	return &Controller{pipeline: pipeline}
}

func (c *Controller) OneTimeSyncCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
}

func (c *Controller) DailySyncCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
}
