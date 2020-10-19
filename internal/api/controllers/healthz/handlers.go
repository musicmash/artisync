package healthz

import (
	"net/http"

	"github.com/musicmash/artisync/internal/db"
)

type Controller struct {
	conn *db.Conn
}

func New(conn *db.Conn) *Controller {
	return &Controller{conn: conn}
}

func (c *Controller) Get(w http.ResponseWriter, _ *http.Request) {
	if err := c.conn.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return //nolint:nlreturn
	}

	w.WriteHeader(http.StatusOK)
}
