package spotify

import (
	"errors"
	"net/url"

	"github.com/musicmash/artisync/internal/log"
)

const (
	stateBackgroundSyncAllowed = "sync-in-background-allowed"
	stateBackgroundSyncDenied  = "sync-in-background-denied"
)

var (
	errAuthFailed   = errors.New("auth failed")
	errCodeIsEmpty  = errors.New("query arg 'code' is empty")
	errUnknownState = errors.New("unknown state")
)

func validateQuery(values url.Values) error {
	if e := values.Get("error"); e != "" {
		log.Errorf("auth failed: %v", e)

		return errAuthFailed
	}

	state := values.Get("state")
	if state != stateBackgroundSyncAllowed && state != stateBackgroundSyncDenied {
		return errUnknownState
	}

	code := values.Get("code")
	if code == "" {
		return errCodeIsEmpty
	}

	return nil
}
