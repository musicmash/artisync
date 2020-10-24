package spotify

import (
	"errors"
	"net/url"
)

const (
	stateBackgroundSyncAllowed = "sync-in-background-allowed"
	stateBackgroundSyncDenied  = "sync-in-background-denied"
)

var (
	errCodeIsEmpty  = errors.New("query arg 'code' is empty")
	errUnknownState = errors.New("unknown state")
)

func validateStateAndCode(values url.Values) error {
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
