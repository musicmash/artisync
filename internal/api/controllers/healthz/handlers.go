package healthz

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/musicmash/artisync/internal/db"
	"github.com/musicmash/artisync/internal/db/models"
)

type Controller struct {
	conn *db.Conn
}

func New(conn *db.Conn) *Controller {
	return &Controller{conn: conn}
}

func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	err := c.conn.ExecTx(context.Background(), func(querier *models.Queries) error {
		art, err := querier.CreateArtist(r.Context(), models.CreateArtistParams{
			Name:   "artisync-test",
			Poster: sql.NullString{},
		})
		if err != nil {
			return fmt.Errorf("can't create new artist: %w", err)
		}

		_, err = querier.CreateArtistAssociation(r.Context(), models.CreateArtistAssociationParams{
			ArtistID:  art.ID,
			StoreName: "spotify",
			StoreID:   "059c3940-a791-422d-8330-2954918c51e6",
		})
		if err != nil {
			return fmt.Errorf("can't associate artist: %w", err)
		}

		err = querier.CreateSubscription(r.Context(), models.CreateSubscriptionParams{
			UserName:  "objque",
			StoreName: "spotify",
			StoreID:   "059c3940-a791-422d-8330-2954918c51e6",
		})
		if err != nil {
			return fmt.Errorf("can't subscribe user: %w", err)
		}

		return nil
	})
	if err != nil {
		_, _ = fmt.Fprintf(w, "got error after tx: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Get(w http.ResponseWriter, _ *http.Request) {
	if err := c.conn.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
